package datasource

import (
	"time"

	"github.com/google/uuid"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/database"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/random"
)

type IServiceRepo interface {
	FindOneById(uuid.UUID) (*Service, error)
	FindOneByName(string) (*Service, error)
	FulltextSearchServiceByName(string) ([]*Service, error)
	FulltextSearchServiceByNameOrInGroup(string, ...string) ([]*Service, error)
	UpdateOneById(uuid.UUID, map[string]interface{}) (*Service, error)
	CreateService(map[string]interface{}) (*Service, error)
	DeleteOneById(uuid.UUID) error
}

type ServiceRepo struct {
	db *database.PostgreDb
}

func (sr *ServiceRepo) FindOneById(id uuid.UUID) (*Service, error) {
	service := &Service{}
	result := sr.db.First(service, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return service, nil
}

func (sr *ServiceRepo) FindOneByName(name string) (*Service, error) {
	service := &Service{}
	result := sr.db.First(service, "name = ?", name)
	if result.Error != nil {
		return nil, result.Error
	}
	return service, nil
}

func (sr *ServiceRepo) UpdateOneById(id uuid.UUID, updateParams map[string]interface{}) (*Service, error) {
	_, err := sr.FindOneById(id)
	if err != nil {
		return nil, err
	}

	result := sr.db.Model(&Service{}).Where("id = ?", id).Updates(updateParams)
	if result.Error != nil {
		return nil, result.Error
	}
	return sr.FindOneById(id)
}

func (sr *ServiceRepo) CreateService(data map[string]interface{}) (*Service, error) {
	var err error

	data["created_at"] = time.Now()
	data["id"] = random.GenerateRandomUUID()

	result := sr.db.Model(&Service{}).Create(data)
	if result.Error != nil {
		return nil, result.Error
	}

	service, err := sr.FindOneById(uuid.MustParse(data["id"].(string)))
	if err != nil {
		return nil, err
	}

	return service, nil
}

func (sr *ServiceRepo) DeleteOneById(id uuid.UUID) error {
	result := sr.db.Delete(&Service{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (sr *ServiceRepo) FulltextSearchServiceByName(name string) ([]*Service, error) {
	// Build the query with full-text search on Name
	query := sr.db.Debug().Where("to_tsvector('simple', name) @@ plainto_tsvector('simple', ?)", name)

	// Apply soft delete filter (if applicable)
	if sr.db.Dialector.Name() == "postgres" {
		query = query.Where("deleted_at IS NULL")
	}

	// Execute the query and scan the results
	var services []*Service
	err := query.Find(&services)
	if err.Error != nil {
		return nil, err.Error
	}
	// Return the results and any errors
	return services, nil
}

func (sr *ServiceRepo) FulltextSearchServiceByNameOrInGroup(name string, groupIds ...string) ([]*Service, error) {
	// Build the query with full-text search on Name
	query := sr.db.Debug().Where("to_tsvector('simple', name) @@ plainto_tsvector('simple', ?)", name).
		Or("group_service_id IN (?)", groupIds)

	// Apply soft delete filter (if applicable)
	if sr.db.Dialector.Name() == "postgres" {
		query = query.Where("deleted_at IS NULL")
	}

	// Execute the query and scan the results
	var services []*Service
	err := query.Find(&services)
	if err.Error != nil {
		return nil, err.Error
	}
	// Return the results and any errors
	return services, nil
}

func NewServiceRepo(db *database.PostgreDb) *ServiceRepo {
	return &ServiceRepo{
		db: db,
	}
}
