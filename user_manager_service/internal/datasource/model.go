package datasource

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID      `json:"id,omitempty"         gorm:"type:uuid;default:uuid_generate_v4();primarykey;index"`
	CreatedAt time.Time      `json:"created_at,omitempty"`
	UpdatedAt time.Time      `json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
	Email     string         `json:"email,omitempty"      gorm:"index"`
	Password  string         `json:"password,omitempty"`
	Phone     string         `json:"phone,omitempty"`
	FirstName string         `json:"first_name,omitempty"`
	LastName  string         `json:"last_name,omitempty"`
	Status    string         `json:"status,omitempty"`
	Cards     []Card
}

type BannedAccount struct {
	ID          uuid.UUID      `json:"id,omitempty"         gorm:"type:uuid;default:uuid_generate_v4();primarykey;index"`
	CreatedAt   time.Time      `json:"created_at,omitempty"`
	UpdatedAt   time.Time      `json:"updated_at,omitempty"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
	UserId      uuid.UUID      `json:"user_id,omitempty"    gorm:"type:uuid,index"`
	Description string         `json:"reason,omitempty"`
	User        User
}

type Card struct {
	ID         uuid.UUID      `json:"id,omitempty"          gorm:"type:uuid;default:uuid_generate_v4();primarykey"`
	CreatedAt  time.Time      `json:"created_at,omitempty"`
	UpdatedAt  time.Time      `json:"updated_at,omitempty"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at,omitempty"  gorm:"index"`
	CardNumber string         `json:"card_number,omitempty"`
	OwnerName  string         `json:"owner_name,omitempty"`
	BankName   string         `json:"bank_name,omitempty"`
	UserId     uuid.UUID      `json:"user_id,omitempty"     gorm:"type:uuid"`
	User       User
}
