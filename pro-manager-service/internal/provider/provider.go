package provider

import (
	proapi "github.com/nguyentrunghieu15/be-beehome-prj/api/pro-api"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/logwrapper"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/validator"
	"github.com/nguyentrunghieu15/be-beehome-prj/pro-manager-service/internal/datasource"
)

type ProviderService struct {
	proapi.ProServiceServer
	logger          logwrapper.ILoggerWrapper
	validator       validator.IValidator
	proRepo         datasource.IProviderRepo
	paymentRepo     datasource.IPaymentMethodRepo
	hireRepo        datasource.IHireRepo
	reviewRepo      datasource.IReviewRepo
	socialMediaRepo datasource.ISocialMediaRepo
	postalCodeRepo  datasource.IPostalCodeRepo
}

type ProviderServiceBuilder struct {
	logger          logwrapper.ILoggerWrapper
	validator       validator.IValidator
	proRepo         datasource.IProviderRepo
	paymentRepo     datasource.IPaymentMethodRepo
	hireRepo        datasource.IHireRepo
	reviewRepo      datasource.IReviewRepo
	socialMediaRepo datasource.ISocialMediaRepo
	postalCodeRepo  datasource.IPostalCodeRepo
}

func (b *ProviderServiceBuilder) SetLogger(l logwrapper.ILoggerWrapper) *ProviderServiceBuilder {
	b.logger = l
	return b
}

func (b *ProviderServiceBuilder) SetValidator(v validator.IValidator) *ProviderServiceBuilder {
	SetRules(v)
	b.validator = v
	return b
}

func (b *ProviderServiceBuilder) SetProRepo(r datasource.IProviderRepo) *ProviderServiceBuilder {
	b.proRepo = r
	return b
}

func (b *ProviderServiceBuilder) SetPaymentMethodRepo(r datasource.IPaymentMethodRepo) *ProviderServiceBuilder {
	b.paymentRepo = r
	return b
}

func (b *ProviderServiceBuilder) SetHireRepo(r datasource.IHireRepo) *ProviderServiceBuilder {
	b.hireRepo = r
	return b
}

func (b *ProviderServiceBuilder) SetReviewRepo(r datasource.IReviewRepo) *ProviderServiceBuilder {
	b.reviewRepo = r
	return b
}

func (b *ProviderServiceBuilder) SetSocialMediaRepo(r datasource.ISocialMediaRepo) *ProviderServiceBuilder {
	b.socialMediaRepo = r
	return b
}

func (b *ProviderServiceBuilder) SetPostalCodeRepo(r datasource.IPostalCodeRepo) *ProviderServiceBuilder {
	b.postalCodeRepo = r
	return b
}

// Build function to create the ProviderService instance
func (b *ProviderServiceBuilder) Build() (*ProviderService, error) {
	// Validate required fields (optional)
	// ...

	return &ProviderService{
		logger:          b.logger,
		validator:       b.validator,
		proRepo:         b.proRepo,
		paymentRepo:     b.paymentRepo,
		hireRepo:        b.hireRepo,
		reviewRepo:      b.reviewRepo,
		socialMediaRepo: b.socialMediaRepo,
		postalCodeRepo:  b.postalCodeRepo,
	}, nil
}
