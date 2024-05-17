package migration

import (
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/database"
	"github.com/nguyentrunghieu15/be-beehome-prj/user_manager_service/internal/datasource"
)

func MigrationDatasource(db *database.PostgreDb) error {
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	return db.AutoMigrate(
		&datasource.User{},
		&datasource.Card{})
}
