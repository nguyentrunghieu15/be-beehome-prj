package hireservice

import (
	"errors"

	proapi "github.com/nguyentrunghieu15/be-beehome-prj/api/pro-api"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/logwrapper"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/validator"
	addressclient "github.com/nguyentrunghieu15/be-beehome-prj/pro-manager-service/internal/address-client"
	"github.com/nguyentrunghieu15/be-beehome-prj/pro-manager-service/internal/datasource"
)

type HireService struct {
	proapi.HireServiceServer
	hireRepo      datasource.IHireRepo
	proRepo       datasource.IProviderRepo
	logger        logwrapper.ILoggerWrapper
	validator     validator.IValidator
	addressClient addressclient.IAddressClient
}

type HireServiceBuilder struct {
	hireRepo      datasource.IHireRepo
	proRepo       datasource.IProviderRepo
	logger        logwrapper.ILoggerWrapper
	validator     validator.IValidator
	addressClient addressclient.IAddressClient
}

// NewHireServiceBuilder creates a new builder instance
func NewHireServiceBuilder() *HireServiceBuilder {
	return &HireServiceBuilder{}
}

// WithHireRepo sets the HireRepo for the service
func (b *HireServiceBuilder) WithHireRepo(repo datasource.IHireRepo) *HireServiceBuilder {
	b.hireRepo = repo
	return b
}

// WithProviderRepo sets the ProviderRepo for the service
func (b *HireServiceBuilder) WithProviderRepo(repo datasource.IProviderRepo) *HireServiceBuilder {
	b.proRepo = repo
	return b
}

// WithLogger sets the logger for the service
func (b *HireServiceBuilder) WithLogger(logger logwrapper.ILoggerWrapper) *HireServiceBuilder {
	b.logger = logger
	return b
}

// WithValidator sets the validator for the service
func (b *HireServiceBuilder) WithValidator(validator validator.IValidator) *HireServiceBuilder {
	b.validator = validator
	SetRules(validator)
	return b
}

func (b *HireServiceBuilder) SetAddressClient(a addressclient.IAddressClient) *HireServiceBuilder {
	b.addressClient = a
	return b
}

// Build constructs the final HireService struct
func (b *HireServiceBuilder) Build() (*HireService, error) {
	// Validate all required fields are set
	if b.hireRepo == nil {
		return nil, errors.New("hireRepo is required")
	}
	if b.proRepo == nil {
		return nil, errors.New("proRepo is required")
	}
	if b.logger == nil {
		return nil, errors.New("logger is required")
	}
	if b.validator == nil {
		return nil, errors.New("validator is required")
	}

	// Build the final HireService struct
	return &HireService{
		hireRepo:      b.hireRepo,
		proRepo:       b.proRepo,
		logger:        b.logger,
		validator:     b.validator,
		addressClient: b.addressClient,
	}, nil
}
