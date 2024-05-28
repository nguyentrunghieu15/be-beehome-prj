package profiles

import (
	"errors"

	userapi "github.com/nguyentrunghieu15/be-beehome-prj/api/user-api"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/logwrapper"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/validator"
	"github.com/nguyentrunghieu15/be-beehome-prj/user-manager-service/internal/datasource"
)

type ProfileService struct {
	userapi.ProfileServiceServer
	validator validator.IValidator
	userRepo  datasource.IUserRepo
	cardRepo  datasource.ICardRepo
	logger    logwrapper.ILoggerWrapper
}

type ProfileServiceBuilder struct {
	validator validator.IValidator
	userRepo  datasource.IUserRepo
	cardRepo  datasource.ICardRepo
	logger    logwrapper.ILoggerWrapper
}

// NewProfileServiceBuilder creates a new ProfileServiceBuilder instance
func NewProfileServiceBuilder() *ProfileServiceBuilder {
	return &ProfileServiceBuilder{}
}

// WithValidator sets the validator for the builder
func (b *ProfileServiceBuilder) WithValidator(v validator.IValidator) *ProfileServiceBuilder {
	SetRules(v)
	b.validator = v
	return b
}

// WithUserRepo sets the user repository for the builder
func (b *ProfileServiceBuilder) WithUserRepo(ur datasource.IUserRepo) *ProfileServiceBuilder {
	b.userRepo = ur
	return b
}

// WithCardRepo sets the card repository for the builder
func (b *ProfileServiceBuilder) WithCardRepo(cr datasource.ICardRepo) *ProfileServiceBuilder {
	b.cardRepo = cr
	return b
}

// WithLogger sets the logger for the builder
func (b *ProfileServiceBuilder) WithLogger(l logwrapper.ILoggerWrapper) *ProfileServiceBuilder {
	b.logger = l
	return b
}

// Build builds and returns a new ProfileService instance
func (b *ProfileServiceBuilder) Build() (*ProfileService, error) {
	if b.validator == nil {
		return nil, errors.New("validator is required")
	}
	if b.userRepo == nil {
		return nil, errors.New("user repository is required")
	}
	if b.cardRepo == nil {
		return nil, errors.New("card repository is required")
	}
	if b.logger == nil {
		return nil, errors.New("logger is required")
	}
	return &ProfileService{
		validator: b.validator,
		userRepo:  b.userRepo,
		cardRepo:  b.cardRepo,
		logger:    b.logger,
	}, nil
}
