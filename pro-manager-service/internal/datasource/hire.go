package datasource

import (
	"time"

	"github.com/google/uuid"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/database"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/random"
)

type IHireRepo interface {
	CreateHire(map[string]interface{}) (*Hire, error)
	FindOneById(id uint) (*Hire, error)
	UpdateHireById(uuid.UUID, map[string]interface{}) (*Hire, error)
	DeleteHire(id uint) error
}

type HireRepo struct {
	db *database.PostgreDb
}

func NewHireRepo(db *database.PostgreDb) *HireRepo {
	return &HireRepo{db: db}
}

func (repo *HireRepo) CreateHire(data map[string]interface{}) (*Hire, error) {

	data["created_at"] = time.Now()
	id := random.GenerateRandomUUID()
	data["id"] = id

	result := repo.db.Model(&Hire{}).Create(data)

	if result.Error != nil {
		return nil, result.Error
	}
	return repo.FindOneById(uuid.MustParse(id))
}

func (repo *HireRepo) FindOneById(id uuid.UUID) (*Hire, error) {
	var hire Hire
	err := repo.db.First(&hire, id).Error
	if err != nil {
		return nil, err
	}
	return &hire, nil
}

func (repo *HireRepo) UpdateHireById(id uuid.UUID, updateParams map[string]interface{}) (*Hire, error) {
	_, err := repo.FindOneById(id)
	if err != nil {
		return nil, err
	}

	result := repo.db.Model(&Hire{}).Where("id = ?", id).Updates(updateParams)
	if result.Error != nil {
		return nil, result.Error
	}
	return repo.FindOneById(id)
}

func (repo *HireRepo) DeleteHire(id uuid.UUID) error {
	var hire Hire
	hire.ID = id
	return repo.db.Delete(&hire).Error
}
