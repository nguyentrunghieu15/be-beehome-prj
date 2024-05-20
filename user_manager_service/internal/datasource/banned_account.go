package datasource

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/database"
	"gorm.io/gorm"
)

type IBannedAccountsRepo interface {
	FindOneByUserId(uuid.UUID) (*BannedAccount, error)
	FindOneById(uuid.UUID) (*BannedAccount, error)
	UpdateOneById(uuid.UUID, map[string]interface{}) (*BannedAccount, error)
	UpdateOneByUserId(uuid.UUID, map[string]interface{}) (*BannedAccount, error)
	CreateBannedAccount(map[string]interface{}) (*BannedAccount, error)
	// DeleteOneById(uuid.UUID) error
}

type BannedAccountsRepo struct {
	db *database.PostgreDb
}

// NewBannedAccountsRepo creates a new instance of BannedAccountsRepo
func NewBannedAccountsRepo(db *database.PostgreDb) *BannedAccountsRepo {
	return &BannedAccountsRepo{db: db}
}

// FindOneByUserId finds a banned account by user ID
func (repo *BannedAccountsRepo) FindOneByUserId(userID uuid.UUID) (*BannedAccount, error) {
	var bannedAccount BannedAccount
	err := repo.db.Where("user_id = ?", userID).First(&bannedAccount).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &bannedAccount, nil
}

// FindOneById finds a banned account by ID
func (repo *BannedAccountsRepo) FindOneById(id uuid.UUID) (*BannedAccount, error) {
	var bannedAccount BannedAccount
	err := repo.db.Where("id = ?", id).First(&bannedAccount).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &bannedAccount, nil
}

// UpdateOneById updates a banned account by ID
func (repo *BannedAccountsRepo) UpdateOneById(
	id uuid.UUID,
	data map[string]interface{},
) (*BannedAccount, error) {
	var bannedAccount BannedAccount
	err := repo.db.Model(&BannedAccount{}).Where("id = ?", id).Updates(data).Error
	if err != nil {
		return nil, err
	}
	err = repo.db.First(&bannedAccount, id).Error
	if err != nil {
		return nil, err
	}
	return &bannedAccount, nil
}

// UpdateOneByUserId updates a banned account by user ID
func (repo *BannedAccountsRepo) UpdateOneByUserId(
	userID uuid.UUID,
	data map[string]interface{},
) (*BannedAccount, error) {
	var bannedAccount BannedAccount
	err := repo.db.Model(&bannedAccount).Where("user_id = ?", userID).Updates(data).Error
	if err != nil {
		return nil, err
	}
	err = repo.db.Where("user_id = ?", userID).First(&bannedAccount).Error
	if err != nil {
		return nil, err
	}
	return &bannedAccount, nil
}

// CreateBannedAccount creates a new banned account
func (repo *BannedAccountsRepo) CreateBannedAccount(
	data map[string]interface{},
) (*BannedAccount, error) {
	data["id"] = uuid.New()
	data["created_at"] = time.Now()
	err := repo.db.Model(&BannedAccount{}).Create(data).Error
	if err != nil {
		return nil, err
	}
	return repo.FindOneById(data["id"].(uuid.UUID))
}
