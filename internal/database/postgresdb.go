package datastore

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/nguyentrunghieu15/be-beehome-prj/internal/logwrapper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type PostgreDb gorm.DB

func (p *PostgreDb) Init() *PostgreDb {
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
		return (*PostgreDb)(db)
	}
}

func (p *PostgreDb) UseLogger(l logwrapper.ILoggerWrapper) {
	p.Logger = logger.New(&PostgreWriter{
		loggerwrapper: l,
	}, logger.Config{
		SlowThreshold:             time.Second, // Slow SQL threshold
		LogLevel:                  logger.Info, // Log level
		IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
		ParameterizedQueries:      true,        // Don't include params in the SQL log
		Colorful:                  true,        // Disable color
	})
}

type PostgreWriter struct {
	loggerwrapper logwrapper.ILoggerWrapper
}

func (lw *PostgreWriter) Printf(format string, value ...interface{}) {
	lw.loggerwrapper.Log(logwrapper.NewStandardMsg(
		fmt.Sprintf(format, value...),
	))
}
