package auth

import (
	pb "github.com/nguyentrunghieu15/be-beehome-prj/api/auth-api"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/logwrapper"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/mail"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/validator"
	"github.com/nguyentrunghieu15/be-beehome-prj/pkg/captcha"
	"github.com/nguyentrunghieu15/be-beehome-prj/pkg/jwt"
	"github.com/nguyentrunghieu15/be-beehome-prj/user_manager_service/internal/datasource"
)

type AuthService struct {
	pb.AuthServiceServer
	validator       validator.IValidator
	userRepo        datasource.IUserRepo
	logger          logwrapper.ILoggerWrapper
	jwtGenerator    jwt.IJsonWebTokenizer
	captcharService captcha.ICaptchaService
	sessionStorage  datasource.ISessionStorage
	mailService     mail.IMailBox
}

// AuthServiceBuilder is a builder for AuthService
type AuthServiceBuilder struct {
	validator       validator.IValidator
	userRepo        datasource.IUserRepo
	logger          logwrapper.ILoggerWrapper
	jwtGenerator    jwt.IJsonWebTokenizer
	captcharService captcha.ICaptchaService
	sessionStorage  datasource.ISessionStorage
	mailService     mail.IMailBox
}

// NewAuthServiceBuilder creates a new AuthServiceBuilder instance
func NewAuthServiceBuilder() *AuthServiceBuilder {
	return &AuthServiceBuilder{}
}

// SetValidator sets the validator to be used by AuthService
func (b *AuthServiceBuilder) SetValidator(v validator.IValidator) *AuthServiceBuilder {
	SetRules(v)
	b.validator = v
	return b
}

// SetUserRepository sets the user repository to be used by AuthService
func (b *AuthServiceBuilder) SetUserRepository(ur datasource.IUserRepo) *AuthServiceBuilder {
	b.userRepo = ur
	return b
}

// SetLogger sets the looger to be used by AuthService
func (b *AuthServiceBuilder) SetLogger(l logwrapper.ILoggerWrapper) *AuthServiceBuilder {
	b.logger = l
	return b
}

// Set SetJWTGenerator sets the jwt to be used by AuthService
func (b *AuthServiceBuilder) SetJWTGenerator(j jwt.IJsonWebTokenizer) *AuthServiceBuilder {
	b.jwtGenerator = j
	return b
}

// Set SetCaptchaService sets the captchaService to be used by AuthService
func (b *AuthServiceBuilder) SetCaptchaService(c captcha.ICaptchaService) *AuthServiceBuilder {
	b.captcharService = c
	return b
}

// Set SetSessionStorage sets the sessionStorage to be used by AuthService
func (b *AuthServiceBuilder) SetSessionStorage(s datasource.ISessionStorage) *AuthServiceBuilder {
	b.sessionStorage = s
	return b
}

// Set SetMailService sets the mailService to be used by AuthService
func (b *AuthServiceBuilder) SetMailService(m mail.IMailBox) *AuthServiceBuilder {
	b.mailService = m
	return b
}

// Build builds and returns an AuthService instance
func (b *AuthServiceBuilder) Build() (*AuthService, error) {
	// Validate required dependencies

	// Create and return AuthService instance
	return &AuthService{
		validator:       b.validator,
		userRepo:        b.userRepo,
		logger:          b.logger,
		jwtGenerator:    b.jwtGenerator,
		captcharService: b.captcharService,
		sessionStorage:  b.sessionStorage,
		mailService:     b.mailService,
	}, nil
}
