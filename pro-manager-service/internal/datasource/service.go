package datasource

import (
	"time"

	"github.com/google/uuid"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/database"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/random"
)

type IServiceRepo interface {
	FindServices(interface{}) ([]*Service, error)
	FindOneById(uuid.UUID) (*Service, error)
	FindOneByName(name string) (*Service, error)
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

func (sr *ServiceRepo) FindServices(req interface{}) ([]*Service, error) {
	return nil, nil
}

func NewServiceRepo(db *database.PostgreDb) *ServiceRepo {
	return &ServiceRepo{
		db: db,
	}
}
