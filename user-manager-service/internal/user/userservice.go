package user

import (
	"errors"

	userapi "github.com/nguyentrunghieu15/be-beehome-prj/api/user-api"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/logwrapper"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/validator"
	"github.com/nguyentrunghieu15/be-beehome-prj/user-manager-service/internal/datasource"
)

type UserService struct {
	userapi.UserServiceServer
	userRepo          datasource.IUserRepo
	bannedAccountRepo datasource.IBannedAccountsRepo
	cardRepo          datasource.ICardRepo
	logger            logwrapper.ILoggerWrapper
	validator         validator.IValidator
}

// userServiceBuilder is a concrete implementation of UserServiceBuilder
type UserServiceBuilder struct {
	userRepo          datasource.IUserRepo
	cardRepo          datasource.ICardRepo
	bannedAccountRepo datasource.IBannedAccountsRepo
	logger            logwrapper.ILoggerWrapper
	validator         validator.IValidator
}

// NewUserServiceBuilder creates a new UserServiceBuilder instance
func NewUserServiceBuilder() *UserServiceBuilder {
	return &UserServiceBuilder{}
}

// SetUserRepo sets the datasource.IUserRepo for the UserService
func (b *UserServiceBuilder) SetUserRepo(repo datasource.IUserRepo) *UserServiceBuilder {
	b.userRepo = repo
	return b
}

// SetCardRepo sets the datasource.ICardRepo for the UserService
func (b *UserServiceBuilder) SetCardRepo(repo datasource.ICardRepo) *UserServiceBuilder {
	b.cardRepo = repo
	return b
}

// SetLogger sets the logwrapper.ILoggerWrapper for the UserService
func (b *UserServiceBuilder) SetLogger(logger logwrapper.ILoggerWrapper) *UserServiceBuilder {
	b.logger = logger
	return b
}

// SetValidator sets the validator.IValidator for the UserService
func (b *UserServiceBuilder) SetValidator(validator validator.IValidator) *UserServiceBuilder {
	SetRules(validator)
	b.validator = validator
	return b
}

// SetBannedAccount sets the banned account repository to be used by AuthService
func (b *UserServiceBuilder) SetBannedAccount(bn datasource.IBannedAccountsRepo) *UserServiceBuilder {
	b.bannedAccountRepo = bn
	return b
}

// Build builds and returns a new UserService instance
func (b *UserServiceBuilder) Build() (*UserService, error) {
	if b.userRepo == nil {
		return nil, errors.New("userRepo is required")
	}
	if b.cardRepo == nil {
		return nil, errors.New("cardRepo is required")
	}
	if b.logger == nil {
		return nil, errors.New("logger is required")
	}
	if b.validator == nil {
		return nil, errors.New("validator is required")
	}
	return &UserService{
		userRepo:          b.userRepo,
		cardRepo:          b.cardRepo,
		logger:            b.logger,
		validator:         b.validator,
		bannedAccountRepo: b.bannedAccountRepo,
	}, nil
}
