package datasource

import (
	"time"

	"github.com/google/uuid"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/database"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/random"
)

type IGroupServiceRepo interface {
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
	result := gr.db.First(groupService, "id = ?", id).Preload("Services")
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

func NewGroupServiceRepo(db *database.PostgreDb) *GroupServiceRepo {
	return &GroupServiceRepo{
		db: db,
	}
}
