package datasource

import (
	"time"

	"github.com/google/uuid"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/database"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/random"
)

type IReviewRepo interface {
	FindReviews(interface{}) ([]*Review, error)
	FindOneById(uuid.UUID) (*Review, error)
	FindReviewsByUserId(string) ([]*Review, error)
	FindReviewsByProviderId(string) ([]*Review, error)
	UpdateOneById(uuid.UUID, map[string]interface{}) (*Review, error)
	CreateReview(map[string]interface{}) (*Review, error)
	DeleteOneById(uuid.UUID) error
}

type ReviewRepo struct {
	db *database.PostgreDb
}

func (rr *ReviewRepo) FindOneById(id uuid.UUID) (*Review, error) {
	review := &Review{}
	result := rr.db.Preload("Provider").Preload("Service").First(review, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return review, nil
}

func (rr *ReviewRepo) FindReviewsByUserId(userId string) ([]*Review, error) {
	var reviews []*Review
	result := rr.db.Preload("Provider").Preload("Service").Find(&reviews, "user_id = ?", userId)
	if result.Error != nil {
		return nil, result.Error
	}
	return reviews, nil
}

func (rr *ReviewRepo) FindReviewsByProviderId(providerId string) ([]*Review, error) {
	var reviews []*Review
	result := rr.db.Preload("Provider").Preload("Service").Find(&reviews, "provider_id = ?", providerId)
	if result.Error != nil {
		return nil, result.Error
	}
	return reviews, nil
}

func (rr *ReviewRepo) UpdateOneById(id uuid.UUID, updateParams map[string]interface{}) (*Review, error) {
	_, err := rr.FindOneById(id)
	if err != nil {
		return nil, err
	}

	result := rr.db.Model(&Review{}).Where("id = ?", id).Updates(updateParams)
	if result.Error != nil {
		return nil, result.Error
	}
	return rr.FindOneById(id)
}

func (rr *ReviewRepo) CreateReview(data map[string]interface{}) (*Review, error) {
	var err error

	data["created_at"] = time.Now()
	data["id"] = random.GenerateRandomUUID()

	result := rr.db.Model(&Review{}).Create(data)
	if result.Error != nil {
		return nil, result.Error
	}

	review, err := rr.FindOneById(uuid.MustParse(data["id"].(string)))
	if err != nil {
		return nil, err
	}

	return review, nil
}

func (rr *ReviewRepo) DeleteOneById(id uuid.UUID) error {
	result := rr.db.Delete(&Review{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (rr *ReviewRepo) FindReviews(req interface{}) ([]*Review, error) {
	return nil, nil
}

func NewReviewRepo(db *database.PostgreDb) *ReviewRepo {
	return &ReviewRepo{
		db: db,
	}
}
