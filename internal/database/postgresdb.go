package database

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
		}), &gorm.Config{
			Logger: logger.New(logwrapper.NewLoggerWrapper(), logger.Config{
				LogLevel:                  logger.Info,
				SlowThreshold:             time.Millisecond,
				ParameterizedQueries:      true,
				IgnoreRecordNotFoundError: false,
				Colorful:                  false,
			}),
		})
		if err != nil {
			log.Println(errors.New("postgres error: cant establish connection"))
			log.Printf("postgres error: %v", err)
			time.Sleep(time.Millisecond * 300)
			continue
		}
		return (*PostgreDb)(db)
	}
}
