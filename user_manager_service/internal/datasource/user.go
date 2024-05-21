package datasource

import (
	"errors"
	"time"

	"github.com/google/uuid"
	argon "github.com/nguyentrunghieu15/be-beehome-prj/internal/crypto/argon2"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/database"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/random"
)

type IUserRepo interface {
	FindUsers(interface{}) ([]*User, error)
	FindOneById(uuid.UUID) (*User, error)
	FindOneByEmail(email string) (*User, error)
	UpdateOneById(uuid.UUID, map[string]interface{}) (*User, error)
	CreateUser(map[string]interface{}) (*User, error)
	DeleteOneById(uuid.UUID) error
	FindOneProfileById(uuid.UUID) (*User, error)
}

type UserRepo struct {
	db *database.PostgreDb
}

func (ur *UserRepo) FindOneById(id uuid.UUID) (*User, error) {
	user := &User{}
	result := ur.db.First(user, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (ur *UserRepo) FindOneProfileById(id uuid.UUID) (*User, error) {
	user := &User{}
	result := ur.db.Preload("Cards").First(user, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (ur *UserRepo) FindOneByEmail(email string) (*User, error) {
	user := &User{}
	result := ur.db.First(user, "email = ?", email)
	if result.Error != nil {
		return nil, result.Error
	}
	user.Email = email
	return user, nil
}

func (ur *UserRepo) UpdateOneById(id uuid.UUID, updateParams map[string]interface{}) (*User, error) {
	_, err := ur.FindOneById(id)
	if err != nil {
		return nil, err
	}

	result := ur.db.Model(&User{}).Where("id = ?", id).Updates(updateParams)
	if result.Error != nil {
		return nil, result.Error
	}
	return ur.FindOneById(id)
}

func (ur *UserRepo) CreateUser(data map[string]interface{}) (*User, error) {
	if _, err := ur.FindOneByEmail(data["email"].(string)); err == nil {
		return nil, errors.New("exist email")
	}

	var err error
	if password, ok := data["password"].(string); ok {
		data["password"], err = argon.EncodePassword(password)
		if err != nil {
			return nil, err
		}
	}

	data["created_at"] = time.Now()
	data["id"] = random.GenerateRandomUUID()

	result := ur.db.Model(&User{}).Create(data)
	if result.Error != nil {
		return nil, result.Error
	}

	user, err := ur.FindOneById(uuid.MustParse(data["id"].(string)))
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (ur *UserRepo) DeleteOneById(id uuid.UUID) error {
	result := ur.db.Delete(&User{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (ur *UserRepo) FindUsers(req interface{}) ([]*User, error) {
	return nil, nil
}

func NewUserRepo(db *database.PostgreDb) *UserRepo {
	return &UserRepo{
		db: db,
	}
}
