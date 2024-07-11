package datasource

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	proapi "github.com/nguyentrunghieu15/be-beehome-prj/api/pro-api"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/database"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/random"
)

type IHireRepo interface {
	CreateHire(map[string]interface{}) (*Hire, error)
	FindOneById(uuid.UUID) (*Hire, error)
	FindAll(map[string]interface{}) ([]*Hire, error)
	UpdateHireById(uuid.UUID, map[string]interface{}) (*Hire, error)
	DeleteHire(uuid.UUID) error
	FindByRequest(*proapi.FindHireRequest) ([]*Hire, error)
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
	// Bắt đầu transaction
	tx := repo.db.Begin()
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()

	// Kiểm tra lỗi transaction
	if tx.Error != nil {
		return tx.Error
	}

	// Xóa các review liên quan đến hire
	if err := tx.Where("hire_id = ?", id).Delete(&Review{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Xóa hire
	if err := tx.Where("id = ?", id).Delete(&Hire{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit transaction
	return tx.Commit().Error
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
	if err := db.Preload("Service").
		Preload("Provider").
		Preload("Review").
		Find(&hires).Error; err != nil {
		return nil, err
	}
	return hires, nil
}

func (r *HireRepo) FindByRequest(req *proapi.FindHireRequest) ([]*Hire, error) {
	var hires []*Hire

	db := r.db.DB.Preload("Service").
		Preload("Provider").
		Preload("Review")

	// Apply dataParams based on the provided map

	if req.UserId != nil {
		db = db.Where("user_id = ?", req.UserId)
	}
	if req.ProviderId != nil {
		db = db.Where("provider_id = ?", req.ProviderId)
	}
	if req.Status != nil {
		db = db.Where("status = ?", req.Status)
	}

	if req.ServiceId != nil && *req.ServiceId != "" {
		db = db.Where("service_id = ?", req.ServiceId)
	}

	if req.SearchName != nil && *req.SearchName != "" {
		db = db.Where(`(
  to_tsvector('simple', unaccent(hires.address)) @@ phraseto_tsquery('simple', unaccent(?)) AND
  hires.address IS NOT NULL
) OR (
  to_tsvector('simple', unaccent(hires.issue)) @@ phraseto_tsquery('simple', unaccent(?)) AND
  hires.issue IS NOT NULL
)`, req.SearchName, req.SearchName)
	}

	// Apply pagination based on request parameters
	if req.Pagination != nil {
		// Add sorting logic based on req.Pagination.Sort and req.Pagination.SortBy
		if req.Pagination.Sort != nil && req.Pagination.SortBy != nil {
			sortField := *req.Pagination.SortBy
			sortOrder := "ASC"
			if *req.Pagination.Sort == proapi.TypeSort_DESC {
				sortOrder = "DESC"
			}
			db = db.Order(fmt.Sprintf("%s %s", sortField, sortOrder))
		}

		if req.Pagination.Page != nil && req.Pagination.PageSize != nil {
			offset := int(*req.Pagination.Page) * int(*req.Pagination.PageSize)
			db = db.Offset(offset)
		}

		if req.Pagination.Limit != nil {
			db = db.Limit(int(*req.Pagination.Limit))
		}

	}

	// Find all matching hires
	if err := db.Find(&hires).Error; err != nil {
		return nil, err
	}

	return hires, nil
}
