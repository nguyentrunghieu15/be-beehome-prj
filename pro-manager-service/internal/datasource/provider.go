package datasource

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	proapi "github.com/nguyentrunghieu15/be-beehome-prj/api/pro-api"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/database"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/random"
	"gorm.io/gorm"
)

type IProviderRepo interface {
	FindProviders(*proapi.FindProsRequest) ([]*ProviderReviewInfor, error)
	FindOneById(uuid.UUID) (*Provider, error)
	FindOneByUserId(uuid.UUID) (*Provider, error)
	FindOneByName(name string) (*Provider, error)
	GetAllServicesOfProvider(uuid.UUID) ([]*Service, error)
	UpdateOneById(uuid.UUID, map[string]interface{}) (*Provider, error)
	CreateProvider(map[string]interface{}) (*Provider, error)
	AddServicesForPro(uuid.UUID, ...uuid.UUID) error
	DeleteOneById(uuid.UUID) error
	RemoveServicesOfPro(uuid.UUID, ...uuid.UUID) error
}

type ProviderReviewInfor struct {
	Provider
	HireCount     int64
	AverageRating float64
	ReviewCount   int64
}

type ProviderRepo struct {
	db *database.PostgreDb
}

func (pr *ProviderRepo) FindOneById(id uuid.UUID) (*Provider, error) {
	provider := &Provider{}
	result := pr.db.
		Preload("PaymentMethods").
		Preload("SocialMedias").
		Preload("Hires", "status = ?", "approve").
		First(provider, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return provider, nil
}

func (pr *ProviderRepo) FindOneByUserId(id uuid.UUID) (*Provider, error) {
	provider := &Provider{}
	result := pr.db.
		Preload("PaymentMethods").
		Preload("SocialMedias").
		First(provider, "user_id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return provider, nil
}

func (pr *ProviderRepo) FindOneByName(name string) (*Provider, error) {
	provider := &Provider{}
	result := pr.db.
		Preload("PaymentMethods").
		Preload("SocialMedias").
		First(provider, "name = ?", name)
	if result.Error != nil {
		return nil, result.Error
	}
	return provider, nil
}

func (pr *ProviderRepo) UpdateOneById(id uuid.UUID, updateParams map[string]interface{}) (*Provider, error) {
	_, err := pr.FindOneById(id)
	if err != nil {
		return nil, err
	}

	result := pr.db.Model(&Provider{}).Where("id = ?", id).Updates(updateParams)
	if result.Error != nil {
		return nil, result.Error
	}
	return pr.FindOneById(id)
}

func (pr *ProviderRepo) CreateProvider(data map[string]interface{}) (*Provider, error) {
	data["created_at"] = time.Now()
	data["id"] = random.GenerateRandomUUID()

	result := pr.db.Model(&Provider{}).Create(data)
	if result.Error != nil {
		return nil, result.Error
	}

	provider := &Provider{}

	// Check if the ID exists in the data map
	id := uuid.MustParse(data["id"].(string))

	// Fetch the created record using the retrieved ID
	if err := pr.db.First(provider, id).Error; err != nil {
		return nil, err
	}

	return provider, nil
}

func (pr *ProviderRepo) DeleteOneById(id uuid.UUID) error {
	tx := pr.db.Begin()
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Raw("DELETE provider_service WHERE provider_id = ?", id).Scan(nil); err.Error != nil {
		tx.Rollback()
		return err.Error
	}

	if err := tx.Raw("UPDATE reviews SET deleted_at = NOW() WHERE provider_id = ?", id).Scan(nil); err.Error != nil {
		tx.Rollback()
		return err.Error
	}

	if err := tx.Raw("UPDATE hires SET deleted_at = NOW()  WHERE provider_id = ?", id).Scan(nil); err.Error != nil {
		tx.Rollback()
		return err.Error
	}

	result := tx.Delete(&Provider{}, "id = ?", id)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	if err := tx.Commit(); err.Error != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

func (pr *ProviderRepo) AddServicesForPro(providerID uuid.UUID, serviceIDs ...uuid.UUID) error {
	// Check if provider exists
	if _, err := pr.FindOneById(providerID); err != nil {
		return fmt.Errorf("failed to find provider: %v", err)
	}

	// Create a transaction to ensure data consistency
	tx := pr.db.Begin()
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()

	// Loop through each service ID
	for _, serviceID := range serviceIDs {
		// Check if service exists
		var service Service
		if err := pr.db.Model(&Service{}).Where("id = ?", serviceID).First(&service); err.Error != nil {
			tx.Rollback()
			return fmt.Errorf("failed to find service: %v", err.Error)
		}

		// Create association between provider and service using GORM
		var provider Provider

		if err := tx.First(&provider, providerID); err.Error != nil {
			tx.Rollback()
			return fmt.Errorf("failed to find provider: %v", err)
		}

		err := tx.Model(&provider).Association("Services").Append(&service)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to associate service: %v", err)
		}
	}

	// Commit the transaction if all operations succeed
	if err := tx.Commit(); err.Error != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

func (pr *ProviderRepo) GetAllServicesOfProvider(providerId uuid.UUID) ([]*Service, error) {
	var services []*Service
	var provider Provider
	// Preload services using eager loading to avoid N+1 queries
	result := pr.db.Preload("Services").Where("id = ?", providerId).First(&provider)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // No provider found, return empty slice and nil error
		}
		return nil, result.Error
	}
	services = provider.Services
	return services, nil
}

func (pr *ProviderRepo) RemoveServicesOfPro(providerID uuid.UUID, serviceIDs ...uuid.UUID) error {
	// Check if provider exists
	if _, err := pr.FindOneById(providerID); err != nil {
		return fmt.Errorf("failed to find provider: %v", err)
	}

	// Create a transaction to ensure data consistency
	tx := pr.db.Begin()
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()

	// Loop through each service ID
	for _, serviceID := range serviceIDs {
		// Check if service exists
		var service Service
		if err := pr.db.Model(&Service{}).Where("id = ?", serviceID).First(&service); err.Error != nil {
			tx.Rollback()
			return fmt.Errorf("failed to find service: %v", err.Error)
		}

		// Create association between provider and service using GORM
		var provider Provider

		if err := tx.First(&provider, providerID); err.Error != nil {
			tx.Rollback()
			return fmt.Errorf("failed to find provider: %v", err)
		}

		err := tx.Model(&provider).Association("Services").Delete(&service)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to associate service: %v", err)
		}
	}

	// Commit the transaction if all operations succeed
	if err := tx.Commit(); err.Error != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

func fixAddressToTextSearchQuery(address string) string {

	// Split theo dau phẩy
	a := strings.Replace(address, ",", " ", -1)

	// Loại bỏ dấu cách dư thừa
	a = strings.Join(strings.Fields(a), " ")

	// Split theo dau phẩy
	a = strings.Replace(a, "", " <-> ", -1)

	return a
}

func (pr *ProviderRepo) FindProviders(req *proapi.FindProsRequest) ([]*ProviderReviewInfor, error) {
	// Subquery to filter providers offering 'Oven Cleaning' service
	subquery := pr.db.Select("provider_id").Table("providers").
		Joins("JOIN provider_service ON providers.id = provider_service.provider_id").
		Joins("JOIN services ON services.id = provider_service.service_id").
		Joins("JOIN group_services ON group_services.id = services.group_service_id")
	if req.Filter != nil {
		if req.Filter.ServiceName != nil {
			subquery = subquery.Where(
				"services.name = ? OR group_services.name = ?",
				*req.Filter.ServiceName,
				*req.Filter.ServiceName,
			)
		}
	}

	// Main query with joins, filtering, aggregations, and pagination
	var providers []*ProviderReviewInfor

	query := pr.db.Table("providers"). // Preload Address data
						Select("providers.*, COUNT(hires.id) AS hire_count, AVG(reviews.rating) AS average_rating, COUNT(reviews.id) AS review_count").
						Joins("LEFT JOIN hires ON hires.provider_id = providers.id").
						Joins("LEFT JOIN reviews ON reviews.provider_id = providers.id").
						Where("providers.deleted_at IS NULL")

	if req.Filter != nil {
		if req.Filter.ServiceName != nil {
			// Use subquery for efficient filtering
			query = query.Where("providers.id IN (?)", subquery)
		}
	}
	// Apply filters based on request parameters
	if req.Filter != nil {
		if req.Filter.Name != nil {
			query = query.Where("providers.name like %?%", *req.Filter.Name)
		}
		if req.Filter.Address != nil {
			query = query.Where(
				"to_tsvector('simple', unaccent(providers.address)) @@ phraseto_tsquery('simple', unaccent(?))",
				*req.Filter.Address,
			)
		}
		if req.Filter.Years != nil {
			query = query.Where("providers.years = ?", *req.Filter.Years)
		}
		if req.Filter.Introduction != nil {
			query = query.Where("providers.introduction like %?%", *req.Filter.Introduction)
		}
	}
	query = query.Group("providers.id")

	// Apply pagination based on request parameters
	if req.Pagination != nil {
		// Add sorting logic based on req.Pagination.Sort and req.Pagination.SortBy
		if req.Pagination.Sort != nil && req.Pagination.SortBy != nil {
			sortField := *req.Pagination.SortBy
			sortOrder := "ASC"
			if *req.Pagination.Sort == proapi.TypeSort_DESC {
				sortOrder = "DESC"
			}
			query = query.Order(fmt.Sprintf("%s %s", sortField, sortOrder))
		}

		if req.Pagination.Page != nil && req.Pagination.PageSize != nil {
			offset := int(*req.Pagination.Page) * int(*req.Pagination.PageSize)
			query = query.Offset(offset)
		}

		if req.Pagination.Limit != nil {
			query = query.Limit(int(*req.Pagination.Limit))
		}

	}
	err := query.Scan(&providers)
	if err.Error != nil {
		return nil, err.Error
	}
	return providers, nil
}

func NewProviderRepo(db *database.PostgreDb) *ProviderRepo {
	return &ProviderRepo{
		db: db,
	}
}
