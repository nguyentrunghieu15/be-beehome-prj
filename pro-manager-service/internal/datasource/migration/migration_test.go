package migration

import (
	"fmt"
	"testing"

	"github.com/nguyentrunghieu15/be-beehome-prj/internal/database"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestMigrationDatasource(t *testing.T) {
	dns := fmt.Sprintf("host=localhost user=hiro password=1 dbname=beehome-proservice port=5432 sslmode=disable")
	db, _ := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dns,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	MigrationDatasource(&database.PostgreDb{db})
}
