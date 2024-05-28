package datasource

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/database"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/random"
)

type IProviderRepo interface {
	FindProviders(interface{}) ([]*Provider, error)
	FindOneById(uuid.UUID) (*Provider, error)
	FindOneByName(name string) (*Provider, error)
	UpdateOneById(uuid.UUID, map[string]interface{}) (*Provider, error)
	CreateProvider(map[string]interface{}) (*Provider, error)
	DeleteOneById(uuid.UUID) error
}

type ProviderRepo struct {
	db *database.PostgreDb
}

func (pr *ProviderRepo) FindOneById(id uuid.UUID) (*Provider, error) {
	provider := &Provider{}
	result := pr.db.Preload("PostalCode").Preload("PaymentMethods").Preload("SocialMedias").First(provider, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return provider, nil
}

func (pr *ProviderRepo) FindOneByName(name string) (*Provider, error) {
	provider := &Provider{}
	result := pr.db.Preload("PostalCode").Preload("PaymentMethods").Preload("SocialMedias").First(provider, "name = ?", name)
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
	id, ok := data["ID"].(uuid.UUID)
	if !ok {
		return nil, errors.New("failed to retrieve ID from the created data")
	}

	// Fetch the created record using the retrieved ID
	if err := pr.db.First(provider, id).Error; err != nil {
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

func (pr *ProviderRepo) FindProviders(req interface{}) ([]*Provider, error) {
	return nil, nil
}

func NewProviderRepo(db *database.PostgreDb) *ProviderRepo {
	return &ProviderRepo{
		db: db,
	}
}
