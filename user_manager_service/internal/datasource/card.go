package datasource

import (
	"github.com/google/uuid"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/crypto/aes"
	"gorm.io/gorm"
)

type ICardRepo interface {
	CreateCard(map[string]interface{}) (*Card, error)
	FindOneById(uuid.UUID) (*Card, error)
	DeleteOneById(uuid.UUID) error
	FindAllOfUser(uuid.UUID) ([]*Card, error)
}

type CardRepo struct {
	db *gorm.DB
}

func (c *CardRepo) CreateCard(cardData map[string]interface{}) (*Card, error) {
	cryptedCardNumber, err := aes.Encrypt(cardData["card_number"].(string))
	if err != nil {
		return nil, err
	}
	card := &Card{
		ID:         uuid.New(),
		CardNumber: cryptedCardNumber,
		OwnerName:  cardData["owner_name"].(string),
		BankName:   cardData["bank_name"].(string),
		UserId:     cardData["user_id"].(uuid.UUID),
	}

	result := c.db.Create(card)
	if result.Error != nil {
		return nil, result.Error
	}

	return card, nil
}

func (c *CardRepo) FindOneById(id uuid.UUID) (*Card, error) {
	card := &Card{}
	result := c.db.Where("id = ?", id).First(card)

	if result.Error != nil {
		return nil, result.Error
	}
	return card, nil
}

func (c *CardRepo) DeleteOneById(id uuid.UUID) error {
	result := c.db.Where("id = ?", id).Delete(&Card{})
	return result.Error
}

func (c *CardRepo) FindAllOfUser(userId uuid.UUID) ([]*Card, error) {
	var cards []*Card
	result := c.db.Where("user_id = ?", userId).Find(&cards)
	if result.Error != nil {
		return nil, result.Error
	}
	return cards, nil
}
