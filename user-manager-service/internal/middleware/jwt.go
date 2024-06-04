package middleware

import (
	"github.com/golang-jwt/jwt/v5"

	"github.com/labstack/echo/v4"
)

func WrapperJwtFunc() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Access request headers
			user := c.Get("user").(*jwt.Token)
			if claim, ok := user.Claims.(jwt.MapClaims); ok {
				c.Request().Header.Set("Grpc-Metadata-user-id", claim["ID"].(string))
			}
			// Call the next handler in the chain
			return next(c)
		}
	}
}
