package datasource

import (
	"time"

	"github.com/google/uuid"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/database"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/random"
)

type IPaymentMethodRepo interface {
	FindPaymentMethods(interface{}) ([]*PaymentMethod, error)
	FindOneById(uuid.UUID) (*PaymentMethod, error)
	FindOneByName(name string) (*PaymentMethod, error)
	UpdateOneById(uuid.UUID, map[string]interface{}) (*PaymentMethod, error)
	CreatePaymentMethod(map[string]interface{}) (*PaymentMethod, error)
	DeleteOneById(uuid.UUID) error
}

type PaymentMethodRepo struct {
	db *database.PostgreDb
}

func (pr *PaymentMethodRepo) FindOneById(id uuid.UUID) (*PaymentMethod, error) {
	paymentMethod := &PaymentMethod{}
	result := pr.db.First(paymentMethod, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return paymentMethod, nil
}

func (pr *PaymentMethodRepo) FindOneByName(name string) (*PaymentMethod, error) {
	paymentMethod := &PaymentMethod{}
	result := pr.db.First(paymentMethod, "name = ?", name)
	if result.Error != nil {
		return nil, result.Error
	}
	return paymentMethod, nil
}

func (pr *PaymentMethodRepo) UpdateOneById(id uuid.UUID, updateParams map[string]interface{}) (*PaymentMethod, error) {
	_, err := pr.FindOneById(id)
	if err != nil {
		return nil, err
	}

	result := pr.db.Model(&PaymentMethod{}).Where("id = ?", id).Updates(updateParams)
	if result.Error != nil {
		return nil, result.Error
	}
	return pr.FindOneById(id)
}

func (pr *PaymentMethodRepo) CreatePaymentMethod(data map[string]interface{}) (*PaymentMethod, error) {
	var err error

	data["created_at"] = time.Now()
	data["id"] = random.GenerateRandomUUID()

	result := pr.db.Model(&PaymentMethod{}).Create(data)
	if result.Error != nil {
		return nil, result.Error
	}

	paymentMethod, err := pr.FindOneById(uuid.MustParse(data["id"].(string)))
	if err != nil {
		return nil, err
	}

	return paymentMethod, nil
}

func (pr *PaymentMethodRepo) DeleteOneById(id uuid.UUID) error {
	result := pr.db.Delete(&PaymentMethod{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (pr *PaymentMethodRepo) FindPaymentMethods(req interface{}) ([]*PaymentMethod, error) {
	return nil, nil
}

func NewPaymentMethodRepo(db *database.PostgreDb) *PaymentMethodRepo {
	return &PaymentMethodRepo{
		db: db,
	}
}
