package migration

import "github.com/nguyentrunghieu15/be-beehome-prj/internal/database"

func MigrationDatasource(db *database.PostgreDb) error {
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	return db.AutoMigrate()
}
