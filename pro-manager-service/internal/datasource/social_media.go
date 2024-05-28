package datasource

import (
	"time"

	"github.com/google/uuid"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/database"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/random"
)

type ISocialMediaRepo interface {
	FindSocialMedias(interface{}) ([]*SocialMedia, error)
	FindOneById(uuid.UUID) (*SocialMedia, error)
	FindOneByName(name string) (*SocialMedia, error)
	UpdateOneById(uuid.UUID, map[string]interface{}) (*SocialMedia, error)
	CreateSocialMedia(map[string]interface{}) (*SocialMedia, error)
	DeleteOneById(uuid.UUID) error
}

type SocialMediaRepo struct {
	db *database.PostgreDb
}

func (sr *SocialMediaRepo) FindOneById(id uuid.UUID) (*SocialMedia, error) {
	socialMedia := &SocialMedia{}
	result := sr.db.First(socialMedia, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return socialMedia, nil
}

func (sr *SocialMediaRepo) FindOneByName(name string) (*SocialMedia, error) {
	socialMedia := &SocialMedia{}
	result := sr.db.First(socialMedia, "name = ?", name)
	if result.Error != nil {
		return nil, result.Error
	}
	return socialMedia, nil
}

func (sr *SocialMediaRepo) UpdateOneById(id uuid.UUID, updateParams map[string]interface{}) (*SocialMedia, error) {
	_, err := sr.FindOneById(id)
	if err != nil {
		return nil, err
	}

	result := sr.db.Model(&SocialMedia{}).Where("id = ?", id).Updates(updateParams)
	if result.Error != nil {
		return nil, result.Error
	}
	return sr.FindOneById(id)
}

func (sr *SocialMediaRepo) CreateSocialMedia(data map[string]interface{}) (*SocialMedia, error) {
	var err error

	data["created_at"] = time.Now()
	data["id"] = random.GenerateRandomUUID()

	result := sr.db.Model(&SocialMedia{}).Create(data)
	if result.Error != nil {
		return nil, result.Error
	}

	socialMedia, err := sr.FindOneById(uuid.MustParse(data["id"].(string)))
	if err != nil {
		return nil, err
	}

	return socialMedia, nil
}

func (sr *SocialMediaRepo) DeleteOneById(id uuid.UUID) error {
	result := sr.db.Delete(&SocialMedia{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (sr *SocialMediaRepo) FindSocialMedias(req interface{}) ([]*SocialMedia, error) {
	return nil, nil
}

func NewSocialMediaRepo(db *database.PostgreDb) *SocialMediaRepo {
	return &SocialMediaRepo{
		db: db,
	}
}
