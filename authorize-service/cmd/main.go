package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/nguyentrunghieu15/be-beehome-prj/authorize-service/internal/auth"
	"github.com/nguyentrunghieu15/be-beehome-prj/authorize-service/internal/user"
)

var addr = ":3133"

func main() {
	e := echo.New()
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(echomiddleware.CORSWithConfig(echomiddleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))

	// Routes
	e.POST("/api/v1/auth/forgot-password", auth.ForgotPasswordHandler)
	e.POST("/api/v1/auth/login", auth.LoginHandler)
	e.POST("/api/v1/auth/refresh-token", auth.RefreshTokenHandler)
	e.POST("/api/v1/auth/reset-password", auth.ResetPasswordHandler)
	e.POST("/api/v1/auth/sign-up", auth.SignUpHandler)

	// user
	e.GET("/api/v1/profile", user.GetProfile)
	e.DELETE("/api/v1/profile", user.DeactiveAccount)
	e.POST("/api/v1/profile/add-card", user.AddCard)
	e.POST("/api/v1/profile/change-mail", user.ChangeEmail)
	e.POST("/api/v1/profile/change-name", user.ChangeName)
	e.GET("/api/v1/user", user.ListUsers)
	e.POST("/api/v1/user", user.CreateUser)
	e.POST("/api/v1/user/block/:id", user.BlockUser)
	e.GET("/api/v1/user/:id", user.GetUser)
	e.DELETE("/api/v1/user/:id", user.DeleteUser)
	e.PATCH("/api/v1/user/:id", user.UpdateUser)

	e.Logger.Fatal(e.Start(":3133"))
}
