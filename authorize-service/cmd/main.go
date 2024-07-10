package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/nguyentrunghieu15/be-beehome-prj/authorize-service/internal/auth"
	"github.com/nguyentrunghieu15/be-beehome-prj/authorize-service/internal/cerbosx"
	communication "github.com/nguyentrunghieu15/be-beehome-prj/authorize-service/internal/comunitication"
	"github.com/nguyentrunghieu15/be-beehome-prj/authorize-service/internal/middleware"
	"github.com/nguyentrunghieu15/be-beehome-prj/authorize-service/internal/pro"
	"github.com/nguyentrunghieu15/be-beehome-prj/authorize-service/internal/user"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/envloader"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/kafkax"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/logwrapper"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/mongox"
	singletonmanager "github.com/nguyentrunghieu15/be-beehome-prj/internal/singleton_manager"
)

var addr = ":3133"

const (
	envfile string = "./authorize-service/.env"
	logDir  string = "./authorize-service/logs/authorize-service.log"
)

var rotateWriterConfig = logwrapper.ConfigRollbackWriter{
	MaxAge:     3,
	MaxSize:    10,
	MaxBackups: 3,
	Compress:   true}

func validateEnverionment() error {
	var rules = map[string]interface{}{
		"JWT_SECRET_KEY":         "required",
		"MONGO_USERNAME":         "required",
		"MONGO_PASSWORD":         "required",
		"MONGO_URI":              "required",
		"MONGO_DATABASE":         "required",
		"CERBOS_ADDRESS":         "required",
		"KAFKA_BOOTSTRAP_SERVER": "required",
	}
	return envloader.MustLoad(envfile, rules)
}

func initObject(manager *singletonmanager.SingletonManager) {
	// Create logger for app
	manager.RegisterInstances(&logwrapper.LoggerWrapper{})
	// set output for logger
	logger, _ := (manager.GetInstance(&logwrapper.LoggerWrapper{})).(*logwrapper.LoggerWrapper)
	fileWriter := logwrapper.NewRollbackWriterFile(logDir, rotateWriterConfig)
	out := io.MultiWriter(os.Stdout, fileWriter)
	logger.SetWriter(out)

}

func main() {

	if err := validateEnverionment(); err != nil {
		log.Panic(err)
	}

	// create singleton manager
	manager := singletonmanager.NewSingletonManager()
	// init object
	initObject(manager)

	logger, _ := manager.GetInstance(&logwrapper.LoggerWrapper{}).(*logwrapper.LoggerWrapper)
	e := echo.New()
	e.Use(echomiddleware.RequestLoggerWithConfig(echomiddleware.RequestLoggerConfig{
		LogURI:      true,
		LogStatus:   true,
		LogRemoteIP: true,
		LogProtocol: true,
		LogLatency:  true,
		LogError:    true,
		LogValuesFunc: func(c echo.Context, v echomiddleware.RequestLoggerValues) error {
			logger.Infor(
				fmt.Sprintf("uri:%v status:%v remoteip:%v protocol:%v latency:%v error:%v",
					v.URI,
					v.Status,
					v.RemoteIP,
					v.Protocol,
					v.Latency,
					v.Error,
				),
			)
			return nil
		},
	}))
	// Middleware
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.CORSWithConfig(echomiddleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"*"},
	}))
	e.Use(middleware.AttachProviderFunc())

	mongox.DefaultClient = mongox.NewClientMongoWrapperWithConfig(&mongox.MomgoClientConfig{
		DatabaseName: os.Getenv("MONGO_DATABASE"),
		Username:     os.Getenv("MONGO_USERNAME"),
		Password:     os.Getenv("MONGO_PASSWORD"),
		Address:      os.Getenv("MONGO_URI"),
		Timeout:      time.Minute,
	})
	mongox.DefaultClient.Client()

	cerbosx.DefaultClient = cerbosx.NewCerbosClientWrapperWithConfig(&cerbosx.CerbosClientConfig{
		CerbosAddress: os.Getenv("CERBOS_ADDRESS"),
	})
	cerbosx.DefaultClient.Setup()

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
	e.POST("/api/v1/providers/delete-service:service_id", pro.DeleteServicePro)
	e.DELETE("/api/v1/providers/delete-social-media:id", pro.DeleteSocialMediaPro)
	e.GET("/api/v1/providers/owner/profile", pro.GetProviderProfile)
	e.POST("/api/v1/providers/reply-review", pro.ReplyReviewPro)
	e.POST("/api/v1/providers/review", pro.ReviewPro)
	e.POST("/api/v1/providers/signup", pro.SignUpPro)
	e.PUT("/api/v1/providers/update-social-media", pro.UpdateSocialMediaPro)
	e.GET("/api/v1/providers/:id", pro.FindProById)
	e.DELETE("/api/v1/providers/:id", pro.DeleteProById)
	e.PUT("/api/v1/providers/:id", pro.UpdatePro)
	e.GET("/api/v1/providers/:id/reviews/all", pro.GetAllReviewsOfProvider)
	e.GET("/api/v1/providers/:id/reviews", pro.GetReviewsOfProvider)
	e.GET("/api/v1/providers/:id/services", pro.GetAllServiceOfProvider)

	// ServiceManagerService routes
	e.GET("/api/v1/services", pro.ListServices)
	e.POST("/api/v1/services", pro.CreateService)
	e.GET("/api/v1/services/fulltext/search", pro.FulltextSearchServices)
	e.GET("/api/v1/services/:id", pro.GetService)
	e.DELETE("/api/v1/services/:id", pro.DeleteService)
	e.PATCH("/api/v1/services/:id", pro.UpdateService)

	communication.ProviderResourceKafka = kafkax.NewKafkaClientWrapperWithConfig(
		&kafkax.KafkaClientConfig{
			Topic:            communication.TOPIC_RESOURCE_PROVIDER,
			BooststrapServer: os.Getenv("KAFKA_BOOTSTRAP_SERVER"),
			Protocall:        "tcp",
			MaxBytes:         10e6,
			TimeoutRead:      time.Second,
			TimeoutWrite:     time.Second,
		},
	)

	communication.UserResourceKafka = kafkax.NewKafkaClientWrapperWithConfig(
		&kafkax.KafkaClientConfig{
			Topic:            communication.TOPIC_RESOURCE_USER,
			BooststrapServer: os.Getenv("KAFKA_BOOTSTRAP_SERVER"),
			Protocall:        "tcp",
			MaxBytes:         10e6,
			TimeoutRead:      time.Second,
			TimeoutWrite:     time.Second,
		},
	)

	communication.ServiceResourceKafka = kafkax.NewKafkaClientWrapperWithConfig(
		&kafkax.KafkaClientConfig{
			Topic:            communication.TOPIC_RESOURCE_SERVICE,
			BooststrapServer: os.Getenv("KAFKA_BOOTSTRAP_SERVER"),
			Protocall:        "tcp",
			MaxBytes:         10e6,
			TimeoutRead:      time.Second,
			TimeoutWrite:     time.Second,
		},
	)

	communication.GroupServiceResourceKafka = kafkax.NewKafkaClientWrapperWithConfig(
		&kafkax.KafkaClientConfig{
			Topic:            communication.TOPIC_RESOURCE_GSERVICE,
			BooststrapServer: os.Getenv("KAFKA_BOOTSTRAP_SERVER"),
			Protocall:        "tcp",
			MaxBytes:         10e6,
			TimeoutRead:      time.Second,
			TimeoutWrite:     time.Second,
		},
	)

	communication.HireResourceKafka = kafkax.NewKafkaClientWrapperWithConfig(
		&kafkax.KafkaClientConfig{
			Topic:            communication.TOPIC_RESOURCE_HIRE,
			BooststrapServer: os.Getenv("KAFKA_BOOTSTRAP_SERVER"),
			Protocall:        "tcp",
			MaxBytes:         10e6,
			TimeoutRead:      time.Second,
			TimeoutWrite:     time.Second,
		},
	)

	communication.SocialMediaResourceKafka = kafkax.NewKafkaClientWrapperWithConfig(
		&kafkax.KafkaClientConfig{
			Topic:            communication.TOPIC_RESOURCE_SOCIALMEDIA,
			BooststrapServer: os.Getenv("KAFKA_BOOTSTRAP_SERVER"),
			Protocall:        "tcp",
			MaxBytes:         10e6,
			TimeoutRead:      time.Second,
			TimeoutWrite:     time.Second,
		},
	)

	communication.PaymentMethodResourceKafka = kafkax.NewKafkaClientWrapperWithConfig(
		&kafkax.KafkaClientConfig{
			Topic:            communication.TOPIC_RESOURCE_PAYMENTMETHOD,
			BooststrapServer: os.Getenv("KAFKA_BOOTSTRAP_SERVER"),
			Protocall:        "tcp",
			MaxBytes:         10e6,
			TimeoutRead:      time.Second,
			TimeoutWrite:     time.Second,
		},
	)

	// Provider Resource Handler
	providerMessageHandler := communication.NewProviderResourceHandler(logger)
	communication.ProviderResourceKafka.Reader()
	go func(h *communication.ProviderResourceHandler) {
		for {
			msg, err := communication.ProviderResourceKafka.ReadMessage(context.Background())
			if err != nil {
				logger.Error(err.Error())
				continue
			}
			h.Router(msg)
		}
	}(providerMessageHandler)

	// User Resource Handler
	userMessageHandler := communication.NewUserResourceHandler(logger)
	communication.UserResourceKafka.Reader()
	go func(h *communication.UserResourceHandler) {
		for {
			msg, err := communication.UserResourceKafka.ReadMessage(context.Background())
			if err != nil {
				logger.Error(err.Error())
				continue
			}
			h.Router(msg)
		}
	}(userMessageHandler)

	// Service Resource Handler
	serviceMessageHandler := communication.NewServiceResourceHandler(logger)
	communication.ServiceResourceKafka.Reader()
	go func(h *communication.ServiceResourceHandler) {
		for {
			msg, err := communication.ServiceResourceKafka.ReadMessage(context.Background())
			if err != nil {
				logger.Error(err.Error())
				continue
			}
			h.Router(msg)
		}
	}(serviceMessageHandler)

	// Group Service Resource Handler
	groupServiceMessageHandler := communication.NewGroupServiceResourceHandler(logger)
	communication.GroupServiceResourceKafka.Reader()
	go func(h *communication.GroupServiceResourceHandler) {
		for {
			msg, err := communication.GroupServiceResourceKafka.ReadMessage(context.Background())
			if err != nil {
				logger.Error(err.Error())
				continue
			}
			h.Router(msg)
		}
	}(groupServiceMessageHandler)

	// Hire Resource Handler
	hireMessageHandler := communication.NewHireResourceHandler(logger)
	communication.HireResourceKafka.Reader()
	go func(h *communication.HireResourceHandler) {
		for {
			msg, err := communication.HireResourceKafka.ReadMessage(context.Background())
			if err != nil {
				logger.Error(err.Error())
				continue
			}
			h.Router(msg)
		}
	}(hireMessageHandler)

	// Social Media Resource Handler
	socialMediaMessageHandler := communication.NewSocialMediaResourceHandler(logger)
	communication.SocialMediaResourceKafka.Reader()
	go func(h *communication.SocialMediaResourceHandler) {
		for {
			msg, err := communication.SocialMediaResourceKafka.ReadMessage(context.Background())
			if err != nil {
				logger.Error(err.Error())
				continue
			}
			h.Router(msg)
		}
	}(socialMediaMessageHandler)

	// Payment Method Resource Handler
	paymentMethodMessageHandler := communication.NewPaymentMethodResourceHandler(logger)
	communication.PaymentMethodResourceKafka.Reader()
	go func(h *communication.PaymentMethodResourceHandler) {
		for {
			msg, err := communication.PaymentMethodResourceKafka.ReadMessage(context.Background())
			if err != nil {
				logger.Error(err.Error())
				continue
			}
			h.Router(msg)
		}
	}(paymentMethodMessageHandler)

	defer communication.ProviderResourceKafka.Close()
	defer communication.UserResourceKafka.Close()
	defer communication.ServiceResourceKafka.Close()
	defer communication.GroupServiceResourceKafka.Close()
	defer communication.HireResourceKafka.Close()
	defer communication.SocialMediaResourceKafka.Close()
	defer communication.PaymentMethodResourceKafka.Close()

	e.Logger.Fatal(e.Start(addr))
}
