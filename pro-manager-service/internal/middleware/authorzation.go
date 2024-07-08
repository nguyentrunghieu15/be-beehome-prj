package middleware

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

// AuthMiddleware forwards request to the authorization server with body
func AuthorizationMiddleware(authServiceAddress string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			var copyBody []byte
			if req.Method != "GET" {
				bodyBytes, _ := ioutil.ReadAll(req.Body)
				copyBody = bytes.Clone(bodyBytes)
				req.Body.Close() //  must close
				req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
			}

			// Create a new request to the auth server with the same body
			req, err := http.NewRequest(
				c.Request().Method,
				authServiceAddress+c.Request().RequestURI,
				bytes.NewReader(copyBody),
			)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"message": "failed to create request"})
			}

			// Copy headers from the original request to the new request
			for key, values := range c.Request().Header {
				for _, value := range values {
					req.Header.Add(key, value)
				}
			}

			// Forward the request to the auth server
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				return c.JSON(
					http.StatusInternalServerError,
					map[string]string{"message": "failed to communicate with auth server"},
				)
			}
			defer resp.Body.Close()

			if (resp.StatusCode >= 200 && resp.StatusCode < 300) ||
				strings.HasPrefix(c.Request().RequestURI, "/swagger") {
				return next(c)
			}

			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "failed to unauthorized with server"})
		}
	}
}
