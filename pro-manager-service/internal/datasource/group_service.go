package datasource

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	proapi "github.com/nguyentrunghieu15/be-beehome-prj/api/pro-api"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/database"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/random"
	"gorm.io/gorm"
)

type IGroupServiceRepo interface {
	FindGroupServices(*proapi.ListGroupServicesRequest) ([]*GroupService, error)
	FindOneById(uuid.UUID) (*GroupService, error)
	FindOneByName(name string) (*GroupService, error)
	UpdateOneById(uuid.UUID, map[string]interface{}) (*GroupService, error)
	CreateGroupService(map[string]interface{}) (*GroupService, error)
	DeleteOneById(uuid.UUID) error
	FulltextSearchGroupServiceByName(string) ([]*GroupService, error)
}

type GroupServiceRepo struct {
	db *database.PostgreDb
}

func (gr *GroupServiceRepo) FindOneById(id uuid.UUID) (*GroupService, error) {
	groupService := &GroupService{}
	result := gr.db.Preload("Services").First(groupService, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return groupService, nil
}

func (gr *GroupServiceRepo) FindOneByName(name string) (*GroupService, error) {
	groupService := &GroupService{}
	result := gr.db.First(groupService, "name = ?", name).Preload("Services")
	if result.Error != nil {
		return nil, result.Error
	}
	return groupService, nil
}

func (gr *GroupServiceRepo) UpdateOneById(id uuid.UUID, updateParams map[string]interface{}) (*GroupService, error) {
	_, err := gr.FindOneById(id)
	if err != nil {
		return nil, err
	}

	result := gr.db.Model(&GroupService{}).Where("id = ?", id).Updates(updateParams)
	if result.Error != nil {
		return nil, result.Error
	}
	return gr.FindOneById(id)
}

func (gr *GroupServiceRepo) CreateGroupService(data map[string]interface{}) (*GroupService, error) {
	var err error

	data["created_at"] = time.Now()
	data["id"] = random.GenerateRandomUUID()

	result := gr.db.Model(&GroupService{}).Create(data)
	if result.Error != nil {
		return nil, result.Error
	}

	groupService, err := gr.FindOneById(uuid.MustParse(data["id"].(string)))
	if err != nil {
		return nil, err
	}

	return groupService, nil
}

func (gr *GroupServiceRepo) DeleteOneById(id uuid.UUID) error {
	result := gr.db.Delete(&GroupService{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (gr *GroupServiceRepo) FulltextSearchGroupServiceByName(name string) ([]*GroupService, error) {
	// Build the query with full-text search on Name
	query := gr.db.Debug().Where("to_tsvector(name) @@ phraseto_tsquery(?)", name)

	// Execute the query and scan the results
	var groups []*GroupService
	err := query.Find(&groups)
	if err.Error != nil {
		return nil, err.Error
	}

	// Return the results and any errors
	return groups, nil
}

func (gr *GroupServiceRepo) FindGroupServices(req *proapi.ListGroupServicesRequest) ([]*GroupService, error) {
	var gpServices []*GroupService
	query := gr.db.Model(&GroupService{})

	// Apply filters based on request parameters
	if req.Filter != nil {
		if req.Filter.Name != nil {
			query = query.Where("name like %?%", *req.Filter.Name)
		}

		if req.Filter.CreatedAtFrom != nil {
			from, err := time.Parse(time.RFC3339Nano, *req.Filter.CreatedAtFrom)
			if err != nil {
				return nil, fmt.Errorf("invalid created_at_from: %w", err)
			}
			query = query.Where("created_at >= ?", from)
		}

		if req.Filter.CreatedAtTo != nil {
			to, err := time.Parse(time.RFC3339Nano, *req.Filter.CreatedAtTo)
			if err != nil {
				return nil, fmt.Errorf("invalid created_at_to: %w", err)
			}
			query = query.Where("created_at <= ?", to)
		}

		// Add similar filter logic for UpdatedAtFrom and UpdatedAtTo

		if req.Filter.UpdatedAtFrom != nil {
			from, err := time.Parse(time.RFC3339Nano, *req.Filter.UpdatedAtFrom)
			if err != nil {
				return nil, fmt.Errorf("invalid created_at_from: %w", err)
			}
			query = query.Where("created_at >= ?", from)
		}

		if req.Filter.UpdatedAtTo != nil {
			to, err := time.Parse(time.RFC3339Nano, *req.Filter.UpdatedAtTo)
			if err != nil {
				return nil, fmt.Errorf("invalid created_at_to: %w", err)
			}
			query = query.Where("created_at <= ?", to)
		}

	}

	// Apply pagination based on request parameters
	if req.Pagination != nil {
		if req.Pagination.Limit != nil {
			query = query.Limit(int(*req.Pagination.Limit))
		}

		if req.Pagination.Page != nil && req.Pagination.PageSize != nil {
			offset := int(*req.Pagination.Page) * int(*req.Pagination.PageSize)
			query = query.Offset(offset)
		}

		// Add sorting logic based on req.Pagination.Sort and req.Pagination.SortBy
		if req.Pagination.Sort != nil && req.Pagination.SortBy != nil {
			sortField := *req.Pagination.SortBy
			sortOrder := "ASC"
			if *req.Pagination.Sort == proapi.TypeSort_DESC {
				sortOrder = "DESC"
			}
			query = query.Order(fmt.Sprintf("%s %s", sortField, sortOrder))
		}
	}

	if err := query.Find(&gpServices).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No services found, but not an error
		}
		return nil, fmt.Errorf("error finding services: %w", err)
	}

	return gpServices, nil
}

func NewGroupServiceRepo(db *database.PostgreDb) *GroupServiceRepo {
	return &GroupServiceRepo{
		db: db,
	}
}
