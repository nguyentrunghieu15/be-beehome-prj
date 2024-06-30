package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/nguyentrunghieu15/be-beehome-prj/authorize-service/internal/auth"
	"github.com/nguyentrunghieu15/be-beehome-prj/authorize-service/internal/pro"
	"github.com/nguyentrunghieu15/be-beehome-prj/authorize-service/internal/user"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/envloader"
)

var addr = ":3133"

const (
	envfile string = "./authorize-service/.env"
	logDir  string = "./authorize-service/logs/user-service.log"
)

func validateEnverionment() error {
	var rules = map[string]interface{}{
		"JWT_SECRET_KEY": "required",
		"MONGO_USERNAME": "required",
	}
	return envloader.MustLoad(envfile, rules)
}

func main() {

	if err := validateEnverionment(); err != nil {
		log.Panic(err)
	}

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

	// GroupServiceManager routes
	e.GET("/api/v1/group_services", pro.ListGroupServices)
	e.POST("/api/v1/group_services", pro.CreateGroupService)
	e.GET("/api/v1/group_services/fulltext/search", pro.FulltextSearchGroupServices)
	e.GET("/api/v1/group_services/:id", pro.GetGroupService)
	e.DELETE("/api/v1/group_services/:id", pro.DeleteGroupService)
	e.PATCH("/api/v1/group_services/:id", pro.UpdateGroupService)

	// HireService routes
	e.GET("/api/v1/hires", pro.FindAllHire)
	e.POST("/api/v1/hires", pro.CreateHire)
	e.DELETE("/api/v1/hires/:hireId", pro.DeleteHire)
	e.PATCH("/api/v1/hires/:hireId", pro.UpdateStatusHire)

	// ProService routes
	e.GET("/api/v1/providers", pro.FindPros)
	e.POST("/api/v1/providers", pro.JoinAsProvider)
	e.POST("/api/v1/providers/add-payment-method", pro.AddPaymentMethodPro)
	e.POST("/api/v1/providers/add-service", pro.AddServicePro)
	e.POST("/api/v1/providers/add-social-media", pro.AddSocialMediaPro)
	e.POST("/api/v1/providers/delete-service", pro.DeleteServicePro)
	e.DELETE("/api/v1/providers/delete-social-media", pro.DeleteSocialMediaPro)
	e.GET("/api/v1/providers/owner/profile", pro.GetProviderProfile)
	e.POST("/api/v1/providers/reply-review", pro.ReplyReviewPro)
	e.POST("/api/v1/providers/review", pro.ReviewPro)
	e.POST("/api/v1/providers/signup", pro.SignUpPro)
	e.PUT("/api/v1/providers/update-social-media", pro.UpdateSocialMediaPro)
	e.GET("/api/v1/providers/:id", pro.FindProById)
	e.DELETE("/api/v1/providers/:id", pro.DeleteProById)
	e.PUT("/api/v1/providers/:id", pro.UpdatePro)
	e.GET("/api/v1/providers/:id/reviews", pro.GetAllReviewsOfProvider)
	e.GET("/api/v1/providers/:id/services", pro.GetAllServiceOfProvider)

	// ServiceManagerService routes
	e.GET("/api/v1/services", pro.ListServices)
	e.POST("/api/v1/services", pro.CreateService)
	e.GET("/api/v1/services/fulltext/search", pro.FulltextSearchServices)
	e.GET("/api/v1/services/:id", pro.GetService)
	e.DELETE("/api/v1/services/:id", pro.DeleteService)
	e.PATCH("/api/v1/services/:id", pro.UpdateService)

	e.Logger.Fatal(e.Start(":3133"))
}
