package migration

import (
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/database"
	"github.com/nguyentrunghieu15/be-beehome-prj/pro-manager-service/internal/datasource"
)

func MigrationDatasource(db *database.PostgreDb) error {
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	return db.AutoMigrate(
		&datasource.GroupService{},
		&datasource.Service{},
		&datasource.PaymentMethod{},
		&datasource.Provider{},
		&datasource.Review{},
		&datasource.SocialMedia{},
		&datasource.Hire{},
	)
}
