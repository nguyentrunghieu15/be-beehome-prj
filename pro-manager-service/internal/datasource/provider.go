package datasource

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/database"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/random"
	"gorm.io/gorm"
)

type IProviderRepo interface {
	FindProviders(interface{}) ([]*Provider, error)
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

type ProviderRepo struct {
	db *database.PostgreDb
}

func (pr *ProviderRepo) FindOneById(id uuid.UUID) (*Provider, error) {
	provider := &Provider{}
	result := pr.db.Preload("PostalCode").
		Preload("PaymentMethods").
		Preload("SocialMedias").
		First(provider, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return provider, nil
}

func (pr *ProviderRepo) FindOneByUserId(id uuid.UUID) (*Provider, error) {
	provider := &Provider{}
	result := pr.db.Preload("PostalCode").
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
	result := pr.db.Preload("PostalCode").
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
	if err := pr.db.Preload("PostalCode").First(provider, id).Error; err != nil {
		return nil, err
	}

	return provider, nil
}

func (pr *ProviderRepo) DeleteOneById(id uuid.UUID) error {
	result := pr.db.Delete(&Provider{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
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
			panic(err)
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
			panic(err)
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

func (pr *ProviderRepo) FindProviders(req interface{}) ([]*Provider, error) {
	return nil, nil
}

func NewProviderRepo(db *database.PostgreDb) *ProviderRepo {
	return &ProviderRepo{
		db: db,
	}
}
