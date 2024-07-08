package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	proapi "github.com/nguyentrunghieu15/be-beehome-prj/api/pro-api"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/database"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/envloader"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/kafkax"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/logwrapper"
	singletonmanager "github.com/nguyentrunghieu15/be-beehome-prj/internal/singleton_manager"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/validator"
	"github.com/nguyentrunghieu15/be-beehome-prj/pkg/jwt"
	communication "github.com/nguyentrunghieu15/be-beehome-prj/pro-manager-service/internal/comunitication"
	"github.com/nguyentrunghieu15/be-beehome-prj/pro-manager-service/internal/datasource"
	"github.com/nguyentrunghieu15/be-beehome-prj/pro-manager-service/internal/datasource/migration"
	groupservicemanager "github.com/nguyentrunghieu15/be-beehome-prj/pro-manager-service/internal/groupservice-manager"
	"github.com/nguyentrunghieu15/be-beehome-prj/pro-manager-service/internal/hireservice"
	"github.com/nguyentrunghieu15/be-beehome-prj/pro-manager-service/internal/middleware"
	"github.com/nguyentrunghieu15/be-beehome-prj/pro-manager-service/internal/provider"
	servicemanager "github.com/nguyentrunghieu15/be-beehome-prj/pro-manager-service/internal/service-manager"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	envfile string = "./pro-manager-service/.env"
	logDir  string = "./pro-manager-service/logs/pro-service.log"
)

var rotateWriterConfig = logwrapper.ConfigRollbackWriter{
	MaxAge:     3,
	MaxSize:    10,
	MaxBackups: 3,
	Compress:   true}

func validateEnverionment() error {
	var rules = map[string]interface{}{
		"JWT_SECRET_KEY":         "required",
		"POSTGRES_HOST":          "required",
		"POSTGRES_USER":          "required",
		"POSTGRES_PASSWORD":      "required",
		"POSTGRES_DBNAME":        "required",
		"POSTGRES_PORT":          "required,numeric",
		"POSTGRES_SSLMODE":       "required,oneof=disable enable",
		"CHIPHER_KEY":            "required",
		"AUTHORIZATION_SERVER":   "required",
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

	// create validator for app
	manager.RegisterInstances(&validator.ValidatorStuctMap{})

	// Create connection to postgres
	manager.RegisterInstances(&database.PostgreDb{})
	(manager.GetInstance(&database.PostgreDb{})).(*database.PostgreDb).UseLogger(
		&database.PostgreLoggerDecorater{
			LoggerWrapper: logger,
		},
	)

	// Create connection to redis
	// manager.RegisterInstances(&database.RedisDb{})

	manager.RegisterInstances(&jwt.CustomJWTTokenizer{})

	// manager.RegisterInstances(&captcha.GGRecaptchaService{})

	// manager.RegisterInstances(&mail.MailBox{})
}

const addr = ":3003"

func main() {
	if err := validateEnverionment(); err != nil {
		log.Panic(err)
	}

	// create singleton manager
	manager := singletonmanager.NewSingletonManager()
	// init object
	initObject(manager)

	// Auto migration data
	migration.MigrationDatasource(manager.GetInstance(&database.PostgreDb{}).(*database.PostgreDb))

	lis, err := net.Listen("tcp", ":3002")
	if err != nil {
		log.Panic(err)
	}
	s := grpc.NewServer(grpc.UnaryInterceptor(middleware.UnaryInterceptor))

	proService, err := provider.NewProviderServiceBuilder().
		SetLogger(manager.GetInstance(&logwrapper.LoggerWrapper{}).(*logwrapper.LoggerWrapper)).
		SetHireRepo(datasource.NewHireRepo(manager.GetInstance(&database.PostgreDb{}).(*database.PostgreDb))).
		SetPaymentMethodRepo(datasource.NewPaymentMethodRepo(manager.GetInstance(&database.PostgreDb{}).(*database.PostgreDb))).
		SetPostalCodeRepo(datasource.NewPostalCodeRepo(manager.GetInstance(&database.PostgreDb{}).(*database.PostgreDb))).
		SetProRepo(datasource.NewProviderRepo(manager.GetInstance(&database.PostgreDb{}).(*database.PostgreDb))).
		SetValidator(manager.GetInstance(&validator.ValidatorStuctMap{}).(*validator.ValidatorStuctMap)).
		SetSocialMediaRepo(datasource.NewSocialMediaRepo(manager.GetInstance(&database.PostgreDb{}).(*database.PostgreDb))).
		SetReviewRepo(datasource.NewReviewRepo(manager.GetInstance(&database.PostgreDb{}).(*database.PostgreDb))).
		SetJwtTokenizer(manager.GetInstance(&jwt.CustomJWTTokenizer{}).(*jwt.CustomJWTTokenizer)).
		Build()
	proapi.RegisterProServiceServer(s, proService)

	hireService, err := hireservice.NewHireServiceBuilder().
		WithLogger(manager.GetInstance(&logwrapper.LoggerWrapper{}).(*logwrapper.LoggerWrapper)).
		WithHireRepo(datasource.NewHireRepo(manager.GetInstance(&database.PostgreDb{}).(*database.PostgreDb))).
		WithProviderRepo(datasource.NewProviderRepo(manager.GetInstance(&database.PostgreDb{}).(*database.PostgreDb))).
		WithValidator(manager.GetInstance(&validator.ValidatorStuctMap{}).(*validator.ValidatorStuctMap)).
		Build()
	proapi.RegisterHireServiceServer(s, hireService)

	serviceServer, err := servicemanager.NewServiceManagerServerBuilder().
		WithServiceRepo(datasource.NewServiceRepo(manager.GetInstance(&database.PostgreDb{}).(*database.PostgreDb))).
		// Provide the appropriate parameters
		WithGroupServiceRepo(datasource.NewGroupServiceRepo(manager.GetInstance(&database.PostgreDb{}).(*database.PostgreDb))).

		// Provide the appropriate parameters
		WithLogger(manager.GetInstance(&logwrapper.LoggerWrapper{}).(*logwrapper.LoggerWrapper)).
		WithValidator(manager.GetInstance(&validator.ValidatorStuctMap{}).(*validator.ValidatorStuctMap)).
		Build()
	proapi.RegisterServiceManagerServiceServer(s, serviceServer)

	groupServiceServer, err := groupservicemanager.NewGroupServiceManagerServerBuilder().
		WithServiceRepo(datasource.NewServiceRepo(manager.GetInstance(&database.PostgreDb{}).(*database.PostgreDb))).
		// Provide the appropriate parameters
		WithGroupServiceRepo(datasource.NewGroupServiceRepo(manager.GetInstance(&database.PostgreDb{}).(*database.PostgreDb))).

		// Provide the appropriate parameters
		WithLogger(manager.GetInstance(&logwrapper.LoggerWrapper{}).(*logwrapper.LoggerWrapper)).
		WithValidator(manager.GetInstance(&validator.ValidatorStuctMap{}).(*validator.ValidatorStuctMap)).
		Build()
	proapi.RegisterGroupServiceManagerServer(s, groupServiceServer)

	go s.Serve(lis)

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
	e.Use(middleware.SecureHeaders())
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.CORSWithConfig(echomiddleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))
	e.Use(middleware.AuthorizationMiddleware(os.Getenv("AUTHORIZATION_SERVER")))

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	proMux := runtime.NewServeMux()
	proapi.RegisterProServiceHandlerFromEndpoint(context.Background(), proMux, "localhost:3002", opts)
	e.Any("/api/v1/providers*", echo.WrapHandler(proMux),
		echojwt.WithConfig(echojwt.Config{
			SigningKey: []byte(os.Getenv("JWT_SECRET_KEY")),
			Skipper:    middleware.SkipperJWTProviderService,
		}),
		middleware.WrapperJwtFunc(),
		middleware.AttachProviderFunc(),
	)

	hireMux := runtime.NewServeMux()
	proapi.RegisterHireServiceHandlerFromEndpoint(context.Background(), hireMux, "localhost:3002", opts)
	e.Any("/api/v1/hires*", echo.WrapHandler(hireMux),
		echojwt.WithConfig(echojwt.Config{
			SigningKey: []byte(os.Getenv("JWT_SECRET_KEY")),
		}),
		middleware.WrapperJwtFunc(),
		middleware.AttachProviderFunc(),
	)

	serviceMux := runtime.NewServeMux()
	proapi.RegisterServiceManagerServiceHandlerFromEndpoint(context.Background(), serviceMux, "localhost:3002", opts)
	e.Any("/api/v1/services*", echo.WrapHandler(serviceMux))

	gpServiceMux := runtime.NewServeMux()
	proapi.RegisterGroupServiceManagerHandlerFromEndpoint(context.Background(), gpServiceMux, "localhost:3002", opts)
	e.Any("/api/v1/group_services*", echo.WrapHandler(gpServiceMux))

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

	communication.ReviewResourceKafka = kafkax.NewKafkaClientWrapperWithConfig(
		&kafkax.KafkaClientConfig{
			Topic:            communication.TOPIC_RESOURCE_REVIEW,
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

	defer communication.UserResourceKafka.Close()
	defer communication.HireResourceKafka.Close()
	defer communication.PaymentMethodResourceKafka.Close()
	defer communication.ProviderResourceKafka.Close()
	defer communication.ReviewResourceKafka.Close()
	defer communication.ServiceResourceKafka.Close()
	defer communication.GroupServiceResourceKafka.Close()
	defer communication.SocialMediaResourceKafka.Close()
	communication.HireResourceKafka.Writer()
	communication.PaymentMethodResourceKafka.Writer()
	communication.ProviderResourceKafka.Writer()
	communication.ReviewResourceKafka.Writer()
	communication.ServiceResourceKafka.Writer()
	communication.GroupServiceResourceKafka.Writer()
	communication.SocialMediaResourceKafka.Writer()

	userMessageHandler := communication.NewUserResourceHandler(
		manager.GetInstance(&logwrapper.LoggerWrapper{}).(*logwrapper.LoggerWrapper),
		datasource.NewProviderRepo(manager.GetInstance(&database.PostgreDb{}).(*database.PostgreDb)),
	)

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

	e.Static("/swagger", "./pro-manager-service/static")

	log.Fatal(e.Start(addr))
}
