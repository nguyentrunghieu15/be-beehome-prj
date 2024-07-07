package middleware

import (
	"strings"

	"github.com/labstack/echo/v4"
	customJwt "github.com/nguyentrunghieu15/be-beehome-prj/pkg/jwt"
)

func AttachProviderFunc() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			// Access request headers
			jwtParser := customJwt.CustomJWTTokenizer{}
			authToken := c.Request().Header.Get("Authorization")
			if authToken != "" {
				authToken = strings.Trim(strings.Split(authToken, " ")[1], " ")
				if id, err := jwtParser.ParseToken(authToken); err != nil {
					return err
				} else {
					c.Set("user_id", id)
				}
			}

			providerToken := c.Request().Header.Get("Provider-Id")
			if providerToken != "" {
				if id, err := jwtParser.ParseToken(providerToken); err != nil {
					return err
				} else {
					c.Set("provider_id", id)
				}
			}

			// Call the next handler in the chain
			return next(c)
		}
	}
}
