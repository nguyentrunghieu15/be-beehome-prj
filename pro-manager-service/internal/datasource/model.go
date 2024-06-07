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
	Services  []*Service     `                            gorm:"foreignKey:GroupServiceId"`
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
	Service    Service        `json:"service,omitempty"     gorm:"foreignKey:ServiceId"`
}

type PostalCode struct {
	ID            int32   `json:"id,omitempty"             gorm:"primarykey"`
	CountryCode   string  `json:"country_code,omitempty"`
	Zipcode       string  `json:"zipcode,omitempty"`
	Place         string  `json:"place,omitempty"`
	State         string  `json:"state,omitempty"`
	StateCode     string  `json:"state_code,omitempty"`
	Province      string  `json:"province,omitempty"`
	ProvinceCode  string  `json:"province_code,omitempty"`
	Community     string  `json:"community,omitempty"`
	CommunityCode string  `json:"community_code,omitempty"`
	Latitude      float64 `json:"latitude,omitempty"`
	Longitude     float64 `json:"longitude,omitempty"`
}

type Provider struct {
	ID             uuid.UUID        `json:"id,omitempty"             gorm:"type:uuid;default:uuid_generate_v4();primarykey;index"`
	CreatedAt      time.Time        `json:"created_at,omitempty"`
	CreatedBy      string           `json:"created_by,omitempty"`
	UpdatedAt      time.Time        `json:"updated_at,omitempty"`
	UpdatedBy      string           `json:"updated_by,omitempty"`
	DeletedBy      string           `json:"deleted_by,omitempty"`
	DeletedAt      gorm.DeletedAt   `json:"deleted_at,omitempty"     gorm:"index"`
	Name           string           `json:"name,omitempty"`
	Introduction   string           `json:"introduction,omitempty"`
	Years          int32            `json:"years,omitempty"`
	PostalCodeId   int32            `json:"postal_code_id,omitempty"`
	PostalCode     PostalCode       `json:"postal_code,omitempty"    gorm:"foreignKey:PostalCodeId"`
	UserId         uuid.UUID        `json:"user_id,omitempty"`
	PaymentMethods []*PaymentMethod `                                gorm:"foreignKey:ProviderId"`
	SocialMedias   []*SocialMedia   `                                gorm:"foreignKey:ProviderId"`
	Hires          []*Hire          `                                gorm:"foreignKey:ProviderId"`
	Services       []*Service       `                                gorm:"many2many:provider_service;"`
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
	Provider        Provider       `json:"provider,omitempty"          gorm:"foreignKey:ProviderId"`
	ServiceId       uuid.UUID      `json:"service_id,omitempty"`
	Service         Service        `json:"service,omitempty"           gorm:"foreignKey:ServiceId"`
	WorkTimeFrom    string         `json:"work_time_from,omitempty"`
	WorkTimeTo      string         `json:"work_time_to,omitempty"`
	Status          string         `json:"status,omitempty"` // pendding , aprove , decline ,
	PaymentMethodId uuid.UUID      `json:"payment_method_id,omitempty"`
	PaymentMethod   PaymentMethod  `json:"payment_method,omitempty"    gorm:"foreignKey:PaymentMethodId"`
}
