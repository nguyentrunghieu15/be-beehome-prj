package servicemanager

import (
	"errors"

	proapi "github.com/nguyentrunghieu15/be-beehome-prj/api/pro-api"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/logwrapper"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/validator"
	"github.com/nguyentrunghieu15/be-beehome-prj/pro-manager-service/internal/datasource"
)

type ServiceManagerServer struct {
	proapi.ServiceManagerServiceServer
	serviceRepo      datasource.IServiceRepo
	groupServiceRepo datasource.IGroupServiceRepo
	logger           logwrapper.ILoggerWrapper
	validator        validator.IValidator
}

type ServiceManagerServerBuilder struct {
	serviceRepo      datasource.IServiceRepo
	groupServiceRepo datasource.IGroupServiceRepo
	logger           logwrapper.ILoggerWrapper
	validator        validator.IValidator
}

// NewServiceManagerServerBuilder creates a new builder for ServiceManagerServer.
func NewServiceManagerServerBuilder() *ServiceManagerServerBuilder {
	return &ServiceManagerServerBuilder{}
}

// WithServiceRepo sets the serviceRepo for the builder.
func (b *ServiceManagerServerBuilder) WithServiceRepo(
	serviceRepo datasource.IServiceRepo,
) *ServiceManagerServerBuilder {
	b.serviceRepo = serviceRepo
	return b
}

// WithGroupServiceRepo sets the groupServiceRepo for the builder.
func (b *ServiceManagerServerBuilder) WithGroupServiceRepo(
	groupServiceRepo datasource.IGroupServiceRepo,
) *ServiceManagerServerBuilder {
	b.groupServiceRepo = groupServiceRepo
	return b
}

// WithLogger sets the logger for the builder.
func (b *ServiceManagerServerBuilder) WithLogger(logger logwrapper.ILoggerWrapper) *ServiceManagerServerBuilder {
	b.logger = logger
	return b
}

// WithValidator sets the validator for the builder.
func (b *ServiceManagerServerBuilder) WithValidator(validator validator.IValidator) *ServiceManagerServerBuilder {
	b.validator = validator
	SetRules(validator)
	return b
}

// Build builds a new ServiceManagerServer instance.
func (b *ServiceManagerServerBuilder) Build() (*ServiceManagerServer, error) {
	if b.serviceRepo == nil {
		return nil, errors.New("serviceRepo is required")
	}
	if b.groupServiceRepo == nil {
		return nil, errors.New("groupServiceRepo is required")
	}
	if b.logger == nil {
		return nil, errors.New("logger is required")
	}
	if b.validator == nil {
		return nil, errors.New("validator is required")
	}
	return &ServiceManagerServer{
		serviceRepo:      b.serviceRepo,
		groupServiceRepo: b.groupServiceRepo,
		logger:           b.logger,
		validator:        b.validator,
	}, nil
}
