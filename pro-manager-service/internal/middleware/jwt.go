package middleware

import (
	"strings"

	"github.com/golang-jwt/jwt/v5"

	"github.com/labstack/echo/v4"
	customJwt "github.com/nguyentrunghieu15/be-beehome-prj/pkg/jwt"
)

func WrapperJwtFunc() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Access request headers
			user, ok := c.Get("user").(*jwt.Token)
			if !ok {
				return next(c)
			}
			if claim, ok := user.Claims.(jwt.MapClaims); ok {
				c.Request().Header.Set("Grpc-Metadata-user-id", claim["ID"].(string))
			}
			// Call the next handler in the chain
			return next(c)
		}
	}
}

func AttachProviderFunc() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if skipAttachProviderJwt(c) {
				return next(c)
			}
			// Access request headers
			jwtParser := customJwt.CustomJWTTokenizer{}
			providerToken := c.Request().Header.Get("Provider-Id")
			if id, err := jwtParser.ParseToken(providerToken); err != nil {
				return err
			} else {
				c.Request().Header.Set("Grpc-Metadata-provider-id", id.(string))
			}

			// Call the next handler in the chain
			return next(c)
		}
	}
}

func skipAttachProviderJwt(c echo.Context) bool {
	if strings.HasSuffix(c.Request().URL.Path, "owner/profile") ||
		strings.HasSuffix(c.Request().URL.Path, "add-service") ||
		strings.HasSuffix(c.Request().URL.Path, "add-payment-method") ||
		strings.HasSuffix(c.Request().URL.Path, "add-social-media") ||
		strings.HasPrefix(c.Request().URL.Path, "/api/v1/providers/delete-service/") {
		return false
	}
	return true
}

func SkipperJWTProviderService(c echo.Context) bool {
	if c.Request().Method == "GET" {
		return true
	}
	return false
}
