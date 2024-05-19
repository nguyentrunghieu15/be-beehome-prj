package main

import (
	"context"
	"io"
	"log"
	"net"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	authapi "github.com/nguyentrunghieu15/be-beehome-prj/api/auth-api"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/database"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/envloader"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/logwrapper"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/mail"
	singletonmanager "github.com/nguyentrunghieu15/be-beehome-prj/internal/singleton_manager"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/validator"
	"github.com/nguyentrunghieu15/be-beehome-prj/pkg/captcha"
	"github.com/nguyentrunghieu15/be-beehome-prj/pkg/jwt"
	"github.com/nguyentrunghieu15/be-beehome-prj/user_manager_service/internal/auth"
	"github.com/nguyentrunghieu15/be-beehome-prj/user_manager_service/internal/datasource"
	"github.com/nguyentrunghieu15/be-beehome-prj/user_manager_service/internal/datasource/migration"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	envfile string = "./user_manager_service/.env"
	logDir  string = "./user_manager_service/logs/user-service.log"
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
	s := grpc.NewServer()

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

	go s.Serve(lis)

	if err != nil {
		log.Panic(err)
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err = authapi.RegisterAuthServiceHandlerFromEndpoint(context.Background(), mux, "localhost:3001", opts)

	if err != nil {
		log.Panic(err)
	}

	e.Any("/api/v1/auth/*", echo.WrapHandler(mux))
	e.Static("/swagger", "./user_manager_service/static")

	log.Fatal(e.Start(addr))
}