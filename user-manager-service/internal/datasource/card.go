package datasource

import (
	"github.com/google/uuid"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/crypto/aes"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/database"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/random"
)

type ICardRepo interface {
	CreateCard(map[string]interface{}) (*Card, error)
	FindOneById(uuid.UUID) (*Card, error)
	DeleteOneById(uuid.UUID) error
	FindAllOfUser(uuid.UUID) ([]*Card, error)
}

type CardRepo struct {
	db *database.PostgreDb
}

func (c *CardRepo) CreateCard(cardData map[string]interface{}) (*Card, error) {
	cryptedCardNumber, err := aes.Encrypt(cardData["card_number"].(string))
	if err != nil {
		return nil, err
	}

	cardData["card_number"] = cryptedCardNumber
	cardData["id"] = random.GenerateRandomUUID()

	result := c.db.Model(&Card{}).Create(cardData)
	if result.Error != nil {
		return nil, result.Error
	}

	return c.FindOneById(cardData["id"].(uuid.UUID))
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

func NewCardRepo(db *database.PostgreDb) *CardRepo {
	return &CardRepo{
		db: db,
	}
}
