package mapper

import (
	"time"

	proapi "github.com/nguyentrunghieu15/be-beehome-prj/api/pro-api"
	"github.com/nguyentrunghieu15/be-beehome-prj/pro-manager-service/internal/datasource"
)

func MapProviderToInfo(provider *datasource.Provider) *proapi.ProviderInfo {
	// Handle potential conversion errors (e.g., time format)
	createdAt := provider.CreatedAt.Format(time.RFC3339Nano)
	deletedAt := ""
	updatedAt := ""

	if !provider.UpdatedAt.IsZero() {
		updatedAt = provider.UpdatedAt.Format(time.RFC3339Nano)
	}

	if !provider.DeletedAt.Time.IsZero() {
		deletedAt = provider.DeletedAt.Time.Format(time.RFC3339Nano)
	}
	return &proapi.ProviderInfo{
		Id:            provider.ID.String(),
		CreatedAt:     createdAt, // Assuming time format
		CreatedBy:     provider.CreatedBy,
		UpdatedAt:     updatedAt, // Assuming time format
		UpdatedBy:     provider.UpdatedBy,
		DeletedBy:     provider.DeletedBy,
		DeletedAt:     deletedAt, // Assuming time format (if applicable)
		Name:          provider.Name,
		Introduction:  provider.Introduction,
		Years:         provider.Years,
		Address:       provider.Address,
		NumHires:      int32(len(provider.Hires)), // Assuming Address is populated correctly,
		SocialMedias:  MapToSocialMedias(provider.SocialMedias),
		PaymentMethod: MapToPaymentMethods(provider.PaymentMethods),
	}
}

func MapToProviderViewInfo(info *datasource.ProviderReviewInfor) *proapi.ProviderViewInfor {
	if info == nil {
		return nil
	}
	// Handle potential conversion errors (e.g., time format)
	createdAt := info.CreatedAt.Format(time.RFC3339Nano)
	deletedAt := ""
	updatedAt := ""

	if !info.UpdatedAt.IsZero() {
		updatedAt = info.UpdatedAt.Format(time.RFC3339Nano)
	}

	if !info.DeletedAt.Time.IsZero() {
		deletedAt = info.DeletedAt.Time.Format(time.RFC3339Nano)
	}
	return &proapi.ProviderViewInfor{
		Id:           info.Provider.ID.String(),
		CreatedAt:    createdAt, // Assuming time format for CreatedAt
		CreatedBy:    info.Provider.CreatedBy,
		UpdatedAt:    updatedAt, // Assuming time format for UpdatedAt
		UpdatedBy:    info.Provider.UpdatedBy,
		DeletedBy:    info.Provider.DeletedBy,
		DeletedAt:    deletedAt, // Assuming time format for DeletedAt
		Name:         info.Provider.Name,
		Introduction: info.Provider.Introduction,
		Years:        info.Provider.Years,
		Address:      info.Address,
		NumHires:     int32(info.HireCount),
		Rating: &proapi.ProviderViewInfor_OverviewRating{
			NumRating: int32(info.ReviewCount),
			AvgRating: info.AverageRating,
		},
	}
}

func MapToProviderViewInfos(providers []*datasource.ProviderReviewInfor) []*proapi.ProviderViewInfor {
	proapiProviders := make([]*proapi.ProviderViewInfor, 0, len(providers))
	for _, p := range providers {
		proapiProviders = append(proapiProviders, MapToProviderViewInfo(p))
	}
	return proapiProviders
}

func MapToPaymentMethod(ds *datasource.PaymentMethod) *proapi.PaymentMethod {
	// Handle potential conversion errors (e.g., UUID to string)
	id := ds.ID.String()
	createdAt := ds.CreatedAt.Format(time.RFC3339Nano)
	updatedAt := ""
	if !ds.UpdatedAt.IsZero() {
		updatedAt = ds.UpdatedAt.Format(time.RFC3339Nano)
	}
	deletedAt := ""
	if !ds.DeletedAt.Time.IsZero() {
		deletedAt = ds.DeletedAt.Time.Format(time.RFC3339Nano)
	}

	// Map fields with matching names
	return &proapi.PaymentMethod{
		Id:        id,
		CreatedAt: createdAt,
		CreatedBy: ds.CreatedBy,
		UpdatedAt: updatedAt,
		UpdatedBy: ds.UpdatedBy,
		DeletedBy: ds.DeletedBy,
		DeletedAt: deletedAt,
		Name:      ds.Name,
	}
}

func MapToSocialMedia(ds *datasource.SocialMedia) *proapi.SocialMedia {
	// Handle potential conversion errors (e.g., UUID to string)
	id := ds.ID.String()
	createdAt := ds.CreatedAt.Format(time.RFC3339Nano)
	updatedAt := ""
	if !ds.UpdatedAt.IsZero() {
		updatedAt = ds.UpdatedAt.Format(time.RFC3339Nano)
	}
	deletedAt := ""
	if !ds.DeletedAt.Time.IsZero() {
		deletedAt = ds.DeletedAt.Time.Format(time.RFC3339Nano)
	}

	// Map fields with matching names
	return &proapi.SocialMedia{
		Id:         id,
		CreatedAt:  createdAt,
		CreatedBy:  ds.CreatedBy,
		UpdatedAt:  updatedAt,
		UpdatedBy:  ds.UpdatedBy,
		DeletedBy:  ds.DeletedBy,
		DeletedAt:  deletedAt,
		Name:       ds.Name,
		Link:       ds.Link,
		ProviderId: ds.ProviderId.String(),
	}
}

func MapToPaymentMethods(ds []*datasource.PaymentMethod) []*proapi.PaymentMethod {
	proapiPaymentMethods := make([]*proapi.PaymentMethod, 0, len(ds))
	for _, p := range ds {
		proapiPaymentMethods = append(proapiPaymentMethods, MapToPaymentMethod(p))
	}
	return proapiPaymentMethods
}

func MapToSocialMedias(ds []*datasource.SocialMedia) []*proapi.SocialMedia {
	proapiSocialMedias := make([]*proapi.SocialMedia, 0, len(ds))
	for _, p := range ds {
		proapiSocialMedias = append(proapiSocialMedias, MapToSocialMedia(p))
	}
	return proapiSocialMedias
}
