package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/nguyentrunghieu15/be-beehome-prj/address-service/internal/address"
	grpcaddress "github.com/nguyentrunghieu15/be-beehome-prj/address-service/internal/grpc"
	addressapi "github.com/nguyentrunghieu15/be-beehome-prj/api/address-api"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/database"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/envloader"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/logwrapper"
	singletonmanager "github.com/nguyentrunghieu15/be-beehome-prj/internal/singleton_manager"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/validator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	envfile string = "./address-service/.env"
	logDir  string = "./address-service/logs/user-service.log"
)

var rotateWriterConfig = logwrapper.ConfigRollbackWriter{
	MaxAge:     3,
	MaxSize:    10,
	MaxBackups: 3,
	Compress:   true}

func validateEnverionment() error {
	var rules = map[string]interface{}{
		"POSTGRES_HOST":     "required",
		"POSTGRES_USER":     "required",
		"POSTGRES_PASSWORD": "required",
		"POSTGRES_DBNAME":   "required",
		"POSTGRES_PORT":     "required,numeric",
		"POSTGRES_SSLMODE":  "required,oneof=disable enable",
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

}

const addr = ":3004"

func main() {
	if err := validateEnverionment(); err != nil {
		log.Panic(err)
	}

	// create singleton manager
	manager := singletonmanager.NewSingletonManager()
	// init object
	initObject(manager)

	lis, err := net.Listen("tcp", ":3005")
	if err != nil {
		log.Panic(err)
	}
	s := grpc.NewServer()
	addressHandler := grpcaddress.NewAddressService(address.NewAddressRepo(
		manager.GetInstance(&database.PostgreDb{}).(*database.PostgreDb),
	), manager.GetInstance(&validator.ValidatorStuctMap{}).(*validator.ValidatorStuctMap), manager.GetInstance(&logwrapper.LoggerWrapper{}).(*logwrapper.LoggerWrapper))
	addressapi.RegisterAddressServiceServer(s, addressHandler)
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()

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
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.CORSWithConfig(echomiddleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	mux := runtime.NewServeMux()
	addressapi.RegisterAddressServiceHandlerFromEndpoint(context.Background(), mux, ":3005", opts)
	e.Any("/api/v1/address*", echo.WrapHandler(mux))
	if err != nil {
		log.Panic(err)
	}

	e.Static("/swagger", "./address-service/static")

	log.Fatal(e.Start(addr))
}
