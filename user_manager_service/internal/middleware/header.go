package middleware

import "github.com/labstack/echo/v4"

func SecureHeaders(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("X-Frame-Options", "DENY")
		c.Response().Header().Set("X-XSS-Protection", "1; mode=block")
		c.Response().Header().Set("X-Content-Type-Options", "nosniff")
		c.Response().Header().Set("Referrer-Policy", "same-origin")
		c.Response().Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Response().Header().Set("Content-Security-Policy", "default-src 'self'")
		// Call the next handler
		return next(c)
	}
}
