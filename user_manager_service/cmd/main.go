package main

import (
	"io"
	"log"
	"os"

	"github.com/nguyentrunghieu15/be-beehome-prj/internal/database"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/envloader"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/logwrapper"
	singletonmanager "github.com/nguyentrunghieu15/be-beehome-prj/internal/singleton_manager"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/validator"
	"github.com/nguyentrunghieu15/be-beehome-prj/user_manager_service/internal/datasource/migration"
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
	// Create connection to postgres
	manager.RegisterInstances(&database.PostgreDb{})

	// Create logger for app
	manager.RegisterInstances(&logwrapper.LoggerWrapper{})
	// set output for logger
	fileWriter := logwrapper.NewRollbackWriterFile(logDir, rotateWriterConfig)
	out := io.MultiWriter(os.Stdout, fileWriter)
	(manager.GetInstance(&logwrapper.LoggerWrapper{})).(*logwrapper.LoggerWrapper).SetWriter(out)

	// create validator for app
	manager.RegisterInstances(&validator.ValidatorStuctMap{})
}

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

}
