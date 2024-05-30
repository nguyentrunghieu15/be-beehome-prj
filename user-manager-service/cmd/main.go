package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	authapi "github.com/nguyentrunghieu15/be-beehome-prj/api/auth-api"
	userapi "github.com/nguyentrunghieu15/be-beehome-prj/api/user-api"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/database"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/envloader"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/logwrapper"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/mail"
	singletonmanager "github.com/nguyentrunghieu15/be-beehome-prj/internal/singleton_manager"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/validator"
	"github.com/nguyentrunghieu15/be-beehome-prj/pkg/captcha"
	"github.com/nguyentrunghieu15/be-beehome-prj/pkg/jwt"
	"github.com/nguyentrunghieu15/be-beehome-prj/user-manager-service/internal/auth"
	"github.com/nguyentrunghieu15/be-beehome-prj/user-manager-service/internal/datasource"
	"github.com/nguyentrunghieu15/be-beehome-prj/user-manager-service/internal/datasource/migration"
	"github.com/nguyentrunghieu15/be-beehome-prj/user-manager-service/internal/middleware"
	"github.com/nguyentrunghieu15/be-beehome-prj/user-manager-service/internal/profiles"
	"github.com/nguyentrunghieu15/be-beehome-prj/user-manager-service/internal/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	envfile string = "./user-manager-service/.env"
	logDir  string = "./user-manager-service/logs/user-service.log"
)

var rotateWriterConfig = logwrapper.ConfigRollbackWriter{
	MaxAge:     3,
	MaxSize:    10,
	MaxBackups: 3,
	Compress:   true}

func validateEnverionment() error {
	var rules = map[string]interface{}{
		"JWT_SECRET_KEY":    "required",
		"POSTGRES_HOST":     "required",
		"POSTGRES_USER":     "required",
		"POSTGRES_PASSWORD": "required",
		"POSTGRES_DBNAME":   "required",
		"POSTGRES_PORT":     "required,numeric",
		"POSTGRES_SSLMODE":  "required,oneof=disable enable",
		"CHIPHER_KEY":       "required",
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
	manager.RegisterInstances(&database.RedisDb{})

	manager.RegisterInstances(&jwt.CustomJWTTokenizer{})

	manager.RegisterInstances(&captcha.GGRecaptchaService{})

	manager.RegisterInstances(&mail.MailBox{})
}

const addr = ":3000"

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

	lis, err := net.Listen("tcp", ":3001")
	if err != nil {
		log.Panic(err)
	}
	s := grpc.NewServer(grpc.UnaryInterceptor(middleware.UnaryInterceptor))

	authServer, err := auth.NewAuthServiceBuilder().
		SetLogger(manager.GetInstance(&logwrapper.LoggerWrapper{}).(*logwrapper.LoggerWrapper)).
		SetJWTGenerator(manager.GetInstance(&jwt.CustomJWTTokenizer{}).(*jwt.CustomJWTTokenizer)).
		SetValidator(manager.GetInstance(&validator.ValidatorStuctMap{}).(*validator.ValidatorStuctMap)).
		SetCaptchaService(manager.GetInstance(&captcha.GGRecaptchaService{}).(*captcha.GGRecaptchaService)).
		SetMailService(manager.GetInstance(&mail.MailBox{}).(*mail.MailBox)).
		SetUserRepository(datasource.NewUserRepo(manager.GetInstance(&database.PostgreDb{}).(*database.PostgreDb))).
		SetSessionStorage(datasource.NewSessionStorage(manager.GetInstance(&database.RedisDb{}).(*database.RedisDb))).
		Build()

	authapi.RegisterAuthServiceServer(s, authServer)

	userServer, err := user.NewUserServiceBuilder().
		SetLogger(manager.GetInstance(&logwrapper.LoggerWrapper{}).(*logwrapper.LoggerWrapper)).
		SetValidator(manager.GetInstance(&validator.ValidatorStuctMap{}).(*validator.ValidatorStuctMap)).
		SetUserRepo(datasource.NewUserRepo(manager.GetInstance(&database.PostgreDb{}).(*database.PostgreDb))).
		SetBannedAccount(datasource.NewBannedAccountsRepo(manager.GetInstance(&database.PostgreDb{}).(*database.PostgreDb))).
		SetCardRepo(datasource.NewCardRepo(manager.GetInstance(&database.PostgreDb{}).(*database.PostgreDb))).
		Build()
	if err != nil {
		log.Panic(err)
	}
	userapi.RegisterUserServiceServer(s, userServer)

	profileServer, err := profiles.NewProfileServiceBuilder().
		WithLogger(manager.GetInstance(&logwrapper.LoggerWrapper{}).(*logwrapper.LoggerWrapper)).
		WithValidator(manager.GetInstance(&validator.ValidatorStuctMap{}).(*validator.ValidatorStuctMap)).
		WithUserRepo(datasource.NewUserRepo(manager.GetInstance(&database.PostgreDb{}).(*database.PostgreDb))).
		WithCardRepo(datasource.NewCardRepo(manager.GetInstance(&database.PostgreDb{}).(*database.PostgreDb))).
		Build()
	if err != nil {
		log.Panic(err)
	}
	userapi.RegisterProfileServiceServer(s, profileServer)

	go s.Serve(lis)

	if err != nil {
		log.Panic(err)
	}

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
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	authmux := runtime.NewServeMux()
	authapi.RegisterAuthServiceHandlerFromEndpoint(context.Background(), authmux, "localhost:3001", opts)
	usermux := runtime.NewServeMux()
	userapi.RegisterUserServiceHandlerFromEndpoint(context.Background(), usermux, "localhost:3001", opts)
	userapi.RegisterProfileServiceHandlerFromEndpoint(context.Background(), usermux, "localhost:3001", opts)

	if err != nil {
		log.Panic(err)
	}

	e.Any("/api/v1/auth*", echo.WrapHandler(authmux))
	e.Any("/api/v1/user*", echo.WrapHandler(usermux),
		echojwt.WithConfig(echojwt.Config{
			SigningKey: []byte(os.Getenv("JWT_SECRET_KEY")),
		}),
		middleware.WrapperJwtFunc(),
	)
	e.Any("/api/v1/profile*", echo.WrapHandler(usermux),
		echojwt.WithConfig(echojwt.Config{
			SigningKey: []byte(os.Getenv("JWT_SECRET_KEY")),
		}),
		middleware.WrapperJwtFunc(),
	)
	e.Static("/swagger", "./user-manager-service/static")

	log.Fatal(e.Start(addr))
}
