package datasource

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SocialMedia struct {
	ID         uuid.UUID      `json:"id,omitempty"          gorm:"type:uuid;default:uuid_generate_v4();primarykey;index"`
	CreatedAt  time.Time      `json:"created_at,omitempty"`
	CreatedBy  string         `json:"created_by,omitempty"`
	UpdatedAt  time.Time      `json:"updated_at,omitempty"`
	UpdatedBy  string         `json:"updated_by,omitempty"`
	DeletedBy  string         `json:"deleted_by,omitempty"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at,omitempty"  gorm:"index"`
	Name       string         `json:"name,omitempty"`
	Link       string         `json:"link,omitempty"`
	ProviderId uuid.UUID      `json:"provider_id,omitempty"`
	Provider   Provider       `json:"provider,omitempty"    gorm:"foreignKey:ProviderId"`
}

type PaymentMethod struct {
	ID         uuid.UUID      `json:"id,omitempty"          gorm:"type:uuid;default:uuid_generate_v4();primarykey;index"`
	CreatedAt  time.Time      `json:"created_at,omitempty"`
	CreatedBy  string         `json:"created_by,omitempty"`
	UpdatedAt  time.Time      `json:"updated_at,omitempty"`
	UpdatedBy  string         `json:"updated_by,omitempty"`
	DeletedBy  string         `json:"deleted_by,omitempty"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at,omitempty"  gorm:"index"`
	Name       string         `json:"name,omitempty"`
	ProviderId uuid.UUID      `json:"provider_id,omitempty"`
	Provider   Provider       `json:"provider,omitempty"    gorm:"foreignKey:ProviderId"`
}

type Service struct {
	ID             uuid.UUID      `json:"id,omitempty"         gorm:"type:uuid;default:uuid_generate_v4();primarykey;index"`
	CreatedAt      time.Time      `json:"created_at,omitempty"`
	CreatedBy      string         `json:"created_by,omitempty"`
	UpdatedAt      time.Time      `json:"updated_at,omitempty"`
	UpdatedBy      string         `json:"updated_by,omitempty"`
	DeletedBy      string         `json:"deleted_by,omitempty"`
	DeletedAt      gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
	Name           string         `json:"name,omitempty"`
	Detail         string         `json:"detail,omitempty"`
	Price          float64        `json:"price,omitempty"`
	UnitPrice      string         `json:"unit_price,omitempty"`
	GroupServiceId uuid.UUID
}

type GroupService struct {
	ID        uuid.UUID      `json:"id,omitempty"         gorm:"type:uuid;default:uuid_generate_v4();primarykey;index"`
	CreatedAt time.Time      `json:"created_at,omitempty"`
	CreatedBy string         `json:"created_by,omitempty"`
	UpdatedAt time.Time      `json:"updated_at,omitempty"`
	UpdatedBy string         `json:"updated_by,omitempty"`
	DeletedBy string         `json:"deleted_by,omitempty"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
	Name      string         `json:"name,omitempty"`
	Detail    string         `json:"detail,omitempty"`
	Services  []*Service     `                            gorm:"foreignKey:GroupServiceId;constraint:OnDelete:CASCADE;"`
}

type Review struct {
	ID         uuid.UUID      `json:"id,omitempty"          gorm:"type:uuid;default:uuid_generate_v4();primarykey;index"`
	CreatedAt  time.Time      `json:"created_at,omitempty"`
	CreatedBy  string         `json:"created_by,omitempty"`
	UpdatedAt  time.Time      `json:"updated_at,omitempty"`
	UpdatedBy  string         `json:"updated_by,omitempty"`
	DeletedBy  string         `json:"deleted_by,omitempty"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at,omitempty"  gorm:"index"`
	UserId     uuid.UUID      `json:"user_id,omitempty"`
	ProviderId uuid.UUID      `json:"provider_id,omitempty"`
	Provider   Provider       `json:"provider,omitempty"    gorm:"foreignKey:ProviderId"`
	Rating     int32          `json:"rating,omitempty"`
	Comment    string         `json:"comment,omitempty"`
	Reply      string         `json:"reply,omitempty"`
	ServiceId  uuid.UUID      `json:"service_id,omitempty"`
	Service    *Service       `json:"service,omitempty"     gorm:"foreignKey:ServiceId"`
	Note       *string        `json:"note,omitempty"`
	UserName   string         `json:"user_name,omitempty"`
	HireId     uuid.UUID      `json:"hire_id,omitempty"`
}

type Provider struct {
	ID             uuid.UUID        `json:"id,omitempty"           gorm:"type:uuid;default:uuid_generate_v4();primarykey;index"`
	CreatedAt      time.Time        `json:"created_at,omitempty"`
	CreatedBy      string           `json:"created_by,omitempty"`
	UpdatedAt      time.Time        `json:"updated_at,omitempty"`
	UpdatedBy      string           `json:"updated_by,omitempty"`
	DeletedBy      string           `json:"deleted_by,omitempty"`
	DeletedAt      gorm.DeletedAt   `json:"deleted_at,omitempty"   gorm:"index"`
	Name           string           `json:"name,omitempty"`
	Introduction   string           `json:"introduction,omitempty"`
	Years          int32            `json:"years,omitempty"`
	Address        string           `json:"address,omitempty"`
	UserId         uuid.UUID        `json:"user_id,omitempty"`
	PaymentMethods []*PaymentMethod `                              gorm:"foreignKey:ProviderId;constraint:OnDelete:CASCADE;"`
	SocialMedias   []*SocialMedia   `                              gorm:"foreignKey:ProviderId;constraint:OnDelete:CASCADE;"`
	Hires          []*Hire          `                              gorm:"foreignKey:ProviderId;constraint:OnDelete:CASCADE;"`
	Services       []*Service       `                              gorm:"many2many:provider_service;constraint:OnDelete:CASCADE;"`
}

type Hire struct {
	ID              uuid.UUID      `json:"id,omitempty"                gorm:"type:uuid;default:uuid_generate_v4();primarykey;index"`
	CreatedAt       time.Time      `json:"created_at,omitempty"`
	CreatedBy       string         `json:"created_by,omitempty"`
	UpdatedAt       time.Time      `json:"updated_at,omitempty"`
	UpdatedBy       string         `json:"updated_by,omitempty"`
	DeletedBy       string         `json:"deleted_by,omitempty"`
	DeletedAt       gorm.DeletedAt `json:"deleted_at,omitempty"`
	UserId          string         `json:"user_id,omitempty"`
	ProviderId      uuid.UUID      `json:"provider_id,omitempty"`
	Provider        *Provider      `json:"provider,omitempty"          gorm:"foreignKey:ProviderId"`
	ServiceId       uuid.UUID      `json:"service_id,omitempty"`
	Service         *Service       `json:"service,omitempty"           gorm:"foreignKey:ServiceId"`
	WorkTimeFrom    string         `json:"work_time_from,omitempty"`
	WorkTimeTo      string         `json:"work_time_to,omitempty"`
	Status          string         `json:"status,omitempty"` // pendding , aprove , decline ,
	PaymentMethodId uuid.UUID      `json:"payment_method_id,omitempty"`
	PaymentMethod   *PaymentMethod `json:"payment_method,omitempty"    gorm:"foreignKey:PaymentMethodId"`
	Issue           string         `json:"issue,omitempty"`
	Review          *Review        `json:"review,omitempty"            gorm:"foreignKey:HireId;constraint:OnDelete:CASCADE;"`
	Address         string         `json:"address,omitempty"`
	FullAddress     string         `json:"full_address,omitempty"`
}
