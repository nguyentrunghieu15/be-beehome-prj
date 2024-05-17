package database

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/nguyentrunghieu15/be-beehome-prj/internal/logwrapper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type PostgreDb struct {
	*gorm.DB
}

func (p *PostgreDb) Init() interface{} {
	dns := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=%v",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DBNAME"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_SSLMODE"),
	)
	for {
		db, err := gorm.Open(postgres.New(postgres.Config{
			DSN:                  dns,
			PreferSimpleProtocol: true,
		}), &gorm.Config{})
		if err != nil {
			log.Println(errors.New("postgres error: cant establish connection"))
			log.Printf("postgres error: %v", err)
			time.Sleep(time.Millisecond * 300)
			continue
		}
		return &PostgreDb{DB: db}
	}
}

func (p *PostgreDb) UseLogger(newLogger logger.Writer) {
	p.Config.Logger = logger.New(newLogger, logger.Config{
		SlowThreshold:        time.Millisecond * 300,
		LogLevel:             logger.Info,
		ParameterizedQueries: true,
	})
}

type PostgreLoggerDecorater struct {
	logwrapper.LoggerWrapper
}

func (pl *PostgreLoggerDecorater) Printf(format string, data ...interface{}) {
	format = strings.ReplaceAll(format, "\n", " ")
	pl.LoggerWrapper.Printf(format, data...)
}
