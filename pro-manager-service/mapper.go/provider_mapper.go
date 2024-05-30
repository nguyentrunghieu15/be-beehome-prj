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

	if !provider.UpdatedAt.IsZero() {
		deletedAt = provider.UpdatedAt.Format(time.RFC3339Nano)
	}

	if !provider.DeletedAt.Time.IsZero() {
		deletedAt = provider.DeletedAt.Time.Format(time.RFC3339Nano)
	}
	return &proapi.ProviderInfo{
		Id:           provider.ID.String(),
		CreatedAt:    createdAt, // Assuming time format
		CreatedBy:    provider.CreatedBy,
		UpdatedAt:    provider.UpdatedAt.Format(time.RFC3339Nano), // Assuming time format
		UpdatedBy:    provider.UpdatedBy,
		DeletedBy:    provider.DeletedBy,
		DeletedAt:    deletedAt, // Assuming time format (if applicable)
		Name:         provider.Name,
		Introduction: provider.Introduction,
		Years:        provider.Years,
		PostalCode:   MapPostalCodeToInfo(&provider.PostalCode), // Assuming PostalCode is populated correctly,
	}
}

func MapPostalCodeToInfo(postalCode *datasource.PostalCode) *proapi.PostalCode {
	if postalCode == nil {
		return nil // Handle nil PostalCode gracefully (optional)
	}
	return &proapi.PostalCode{
		Id:            postalCode.ID,
		CountryCode:   postalCode.CountryCode,
		Zipcode:       postalCode.Zipcode,
		Place:         postalCode.Place,
		State:         postalCode.State,
		StateCode:     postalCode.StateCode,
		Province:      postalCode.Province,
		ProvinceCode:  postalCode.ProvinceCode,
		Community:     postalCode.Community,
		CommunityCode: postalCode.CommunityCode,
		Latitude:      postalCode.Latitude,
		Longitude:     postalCode.Longitude,
	}
}
