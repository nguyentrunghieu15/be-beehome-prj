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

type IServiceRepo interface {
	FindServices(*proapi.ListServicesRequest) ([]*Service, error)
	FindOneById(uuid.UUID) (*Service, error)
	FindOneByName(string) (*Service, error)
	FulltextSearchServiceByName(string) ([]*Service, error)
	FulltextSearchServiceByNameOrInGroup(string, ...string) ([]*Service, error)
	UpdateOneById(uuid.UUID, map[string]interface{}) (*Service, error)
	CreateService(map[string]interface{}) (*Service, error)
	DeleteOneById(uuid.UUID) error
}

type ServiceRepo struct {
	db *database.PostgreDb
}

func (sr *ServiceRepo) FindOneById(id uuid.UUID) (*Service, error) {
	service := &Service{}
	result := sr.db.First(service, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return service, nil
}

func (sr *ServiceRepo) FindOneByName(name string) (*Service, error) {
	service := &Service{}
	result := sr.db.First(service, "name = ?", name)
	if result.Error != nil {
		return nil, result.Error
	}
	return service, nil
}

func (sr *ServiceRepo) UpdateOneById(id uuid.UUID, updateParams map[string]interface{}) (*Service, error) {
	_, err := sr.FindOneById(id)
	if err != nil {
		return nil, err
	}

	result := sr.db.Model(&Service{}).Where("id = ?", id).Updates(updateParams)
	if result.Error != nil {
		return nil, result.Error
	}
	return sr.FindOneById(id)
}

func (sr *ServiceRepo) CreateService(data map[string]interface{}) (*Service, error) {
	var err error

	data["created_at"] = time.Now()
	data["id"] = random.GenerateRandomUUID()

	result := sr.db.Model(&Service{}).Create(data)
	if result.Error != nil {
		return nil, result.Error
	}

	service, err := sr.FindOneById(uuid.MustParse(data["id"].(string)))
	if err != nil {
		return nil, err
	}

	return service, nil
}

func (sr *ServiceRepo) DeleteOneById(id uuid.UUID) error {
	result := sr.db.Delete(&Service{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (sr *ServiceRepo) FulltextSearchServiceByName(name string) ([]*Service, error) {
	// Build the query with full-text search on Name
	query := sr.db.Debug().Where("to_tsvector(name) @@ phraseto_tsquery(?)", name)

	// Execute the query and scan the results
	var services []*Service
	err := query.Find(&services)
	if err.Error != nil {
		return nil, err.Error
	}
	// Return the results and any errors
	return services, nil
}

func (sr *ServiceRepo) FulltextSearchServiceByNameOrInGroup(name string, groupIds ...string) ([]*Service, error) {
	// Build the query with full-text search on Name
	query := sr.db.Debug().Where("to_tsvector(name) @@ phraseto_tsquery(?)", name).
		Or("group_service_id IN (?)", groupIds)

	// Execute the query and scan the results
	var services []*Service
	err := query.Find(&services)
	if err.Error != nil {
		return nil, err.Error
	}
	// Return the results and any errors
	return services, nil
}

func (sr *ServiceRepo) FindServices(req *proapi.ListServicesRequest) ([]*Service, error) {
	var services []*Service
	query := sr.db.Model(&Service{})

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

	if err := query.Find(&services).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No services found, but not an error
		}
		return nil, fmt.Errorf("error finding services: %w", err)
	}

	return services, nil
}

func NewServiceRepo(db *database.PostgreDb) *ServiceRepo {
	return &ServiceRepo{
		db: db,
	}
}
