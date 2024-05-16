package datasource

import (
	"time"

	"github.com/google/uuid"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/crypto/aes"
	argon "github.com/nguyentrunghieu15/be-beehome-prj/internal/crypto/argon2"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/random"
	"gorm.io/gorm"
)

type IUserRepo interface {
	FindOneById(uuid.UUID) (*User, error)
	UpdateOneById(uuid.UUID, map[string]interface{}) (*User, error)
	CreateUser(map[string]interface{}) (*User, error)
	DeleteOneById(uuid.UUID) error
}

type UserRepo struct {
	db *gorm.DB
}

func (ur *UserRepo) FindOneById(id uuid.UUID) (*User, error) {
	user := &User{}
	result := ur.db.First(user, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	var err error
	user.Email, err = aes.Decrypt(user.Email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (ur *UserRepo) UpdateOneById(id uuid.UUID, updateParams map[string]interface{}) (*User, error) {
	user, err := ur.FindOneById(id)
	if err != nil {
		return nil, err
	}

	result := ur.db.Model(&User{}).Where("id = ?", id).Updates(updateParams)
	if result.Error != nil {
		return nil, result.Error
	}

	user.Email, err = aes.Decrypt(user.Email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (ur *UserRepo) CreateUser(data map[string]interface{}) (*User, error) {
	var err error
	if email, ok := data["email"].(string); ok {
		data["email"], err = aes.Encrypt(email)
		if err != nil {
			return nil, err
		}
	}
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

	user.Email, err = aes.Decrypt(user.Email)
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
