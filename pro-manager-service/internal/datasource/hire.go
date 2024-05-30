package datasource

import (
	"time"

	"github.com/google/uuid"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/database"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/random"
)

type IHireRepo interface {
	CreateHire(map[string]interface{}) (*Hire, error)
	FindOneById(id uuid.UUID) (*Hire, error)
	FindAll(map[string]interface{}) ([]*Hire, error)
	UpdateHireById(uuid.UUID, map[string]interface{}) (*Hire, error)
	DeleteHire(id uuid.UUID) error
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

func (r *HireRepo) FindAll(dataParams map[string]interface{}) ([]*Hire, error) {
	var hires []*Hire

	db := r.db.DB

	// Apply dataParams based on the provided map
	for field, value := range dataParams {
		switch field {
		case "user_id":
			if userId, ok := value.(string); ok {
				db = db.Where("user_id = ?", userId)
			}
		case "provider_id":
			if providerId, ok := value.(string); ok {
				db = db.Where("provider_id = ?", providerId)
			}
		case "status":
			if status, ok := value.(string); ok {
				db = db.Where("status = ?", status)
			}
		}
	}

	// Find all matching hires
	if err := db.Find(&hires).Error; err != nil {
		return nil, err
	}

	return hires, nil
}
