package datasource

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	proapi "github.com/nguyentrunghieu15/be-beehome-prj/api/pro-api"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/database"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/random"
)

type IReviewRepo interface {
	FindReviews(interface{}) ([]*Review, error)
	FindOneById(uuid.UUID) (*Review, error)
	FindReviewsByUserId(uuid.UUID) ([]*Review, error)
	FindReviewsByProviderId(uuid.UUID) ([]*Review, error)
	FindReviewsByProviderIdWithOptions(uuid.UUID, *proapi.GetReviewOfProviderRequest) ([]*Review, error)
	UpdateOneById(uuid.UUID, map[string]interface{}) (*Review, error)
	CreateReview(map[string]interface{}) (*Review, error)
	DeleteOneById(uuid.UUID) error
}

type ReviewRepo struct {
	db *database.PostgreDb
}

func (rr *ReviewRepo) FindOneById(id uuid.UUID) (*Review, error) {
	review := &Review{}
	result := rr.db.First(review, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return review, nil
}

func (rr *ReviewRepo) FindReviewsByUserId(userId uuid.UUID) ([]*Review, error) {
	var reviews []*Review
	result := rr.db.Find(&reviews, "user_id = ?", userId)
	if result.Error != nil {
		return nil, result.Error
	}
	return reviews, nil
}

func (rr *ReviewRepo) FindReviewsByProviderId(providerId uuid.UUID) ([]*Review, error) {
	var reviews []*Review
	result := rr.db.Preload("Provider").Preload("Service").Find(&reviews, "provider_id = ?", providerId)
	if result.Error != nil {
		return nil, result.Error
	}
	return reviews, nil
}

func (rr *ReviewRepo) FindReviewsByProviderIdWithOptions(
	providerId uuid.UUID,
	req *proapi.GetReviewOfProviderRequest,
) ([]*Review, error) {
	var reviews []*Review
	result := rr.db.Preload("Provider").Preload("Service").Where("provider_id = ?", providerId)
	fmt.Println(req)
	if req != nil && req.Filter != nil {
		if req.Filter.Rating > 0 && req.Filter.Rating < 6 {
			result = result.Where("rating = ?", req.Filter.Rating)
		}
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
			result = result.Order(fmt.Sprintf("%s %s", sortField, sortOrder))
		}

		if req.Pagination.Page != nil && req.Pagination.PageSize != nil {
			offset := int(*req.Pagination.Page) * int(*req.Pagination.PageSize)
			result = result.Offset(offset)
		}

		if req.Pagination.Limit != nil {
			result = result.Limit(int(*req.Pagination.Limit))
		}
	}

	result = result.Find(&reviews)
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
