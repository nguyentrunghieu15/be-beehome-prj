package groupservicemanager

import (
	"errors"

	proapi "github.com/nguyentrunghieu15/be-beehome-prj/api/pro-api"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/logwrapper"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/validator"
	"github.com/nguyentrunghieu15/be-beehome-prj/pro-manager-service/internal/datasource"
)

type GroupServiceManagerServer struct {
	proapi.GroupServiceManagerServer
	serviceRepo      datasource.IServiceRepo
	groupServiceRepo datasource.IGroupServiceRepo
	logger           logwrapper.ILoggerWrapper
	validator        validator.IValidator
}

type GroupServiceManagerServerBuilder struct {
	serviceRepo      datasource.IServiceRepo
	groupServiceRepo datasource.IGroupServiceRepo
	logger           logwrapper.ILoggerWrapper
	validator        validator.IValidator
}

// NewGroupServiceManagerServerBuilder creates a new builder for GroupServiceManagerServer.
func NewGroupServiceManagerServerBuilder() *GroupServiceManagerServerBuilder {
	return &GroupServiceManagerServerBuilder{}
}

// WithServiceRepo sets the serviceRepo for the builder.
func (b *GroupServiceManagerServerBuilder) WithServiceRepo(
	serviceRepo datasource.IServiceRepo,
) *GroupServiceManagerServerBuilder {
	b.serviceRepo = serviceRepo
	return b
}

// WithGroupServiceRepo sets the groupServiceRepo for the builder.
func (b *GroupServiceManagerServerBuilder) WithGroupServiceRepo(
	groupServiceRepo datasource.IGroupServiceRepo,
) *GroupServiceManagerServerBuilder {
	b.groupServiceRepo = groupServiceRepo
	return b
}

// WithLogger sets the logger for the builder.
func (b *GroupServiceManagerServerBuilder) WithLogger(logger logwrapper.ILoggerWrapper) *GroupServiceManagerServerBuilder {
	b.logger = logger
	return b
}

// WithValidator sets the validator for the builder.
func (b *GroupServiceManagerServerBuilder) WithValidator(validator validator.IValidator) *GroupServiceManagerServerBuilder {
	b.validator = validator
	return b
}

// Build builds a new GroupServiceManagerServer instance.
func (b *GroupServiceManagerServerBuilder) Build() (*GroupServiceManagerServer, error) {
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
	return &GroupServiceManagerServer{
		serviceRepo:      b.serviceRepo,
		groupServiceRepo: b.groupServiceRepo,
		logger:           b.logger,
		validator:        b.validator,
	}, nil
}
