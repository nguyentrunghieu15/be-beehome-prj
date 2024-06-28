package mapper

import (
	"time"

	proapi "github.com/nguyentrunghieu15/be-beehome-prj/api/pro-api"
	"github.com/nguyentrunghieu15/be-beehome-prj/pro-manager-service/internal/datasource"
)

func MapToHire(hire *datasource.Hire) *proapi.Hire {
	// Handle potential conversion errors (e.g., time format)
	createdAt := hire.CreatedAt.Format(time.RFC3339Nano)
	deletedAt := ""
	updatedAt := ""

	if !hire.UpdatedAt.IsZero() {
		updatedAt = hire.UpdatedAt.Format(time.RFC3339Nano)
	}

	if !hire.DeletedAt.Time.IsZero() {
		deletedAt = hire.DeletedAt.Time.Format(time.RFC3339Nano)
	}
	return &proapi.Hire{
		Id:              hire.ID.String(),
		CreatedAt:       createdAt,
		CreatedBy:       hire.CreatedBy,
		UpdatedAt:       updatedAt,
		UpdatedBy:       hire.UpdatedBy,
		DeletedBy:       hire.DeletedBy,
		DeletedAt:       deletedAt,
		UserId:          hire.UserId,
		ProviderId:      hire.ProviderId.String(),
		ServiceId:       hire.ServiceId.String(),
		WorkTimeFrom:    hire.WorkTimeFrom,
		WorkTimeTo:      hire.WorkTimeTo,
		Status:          hire.Status,
		PaymentMethodId: hire.PaymentMethodId.String(),
		Issue:           hire.Issue,
	}
}

func MapToListHire(hires []*datasource.Hire) []*proapi.Hire {
	var listHire = make([]*proapi.Hire, 0)
	for _, hire := range hires {
		listHire = append(listHire, MapToHire(hire))
	}
	return listHire
}

func MapToHireInfor(hire *datasource.Hire) *proapi.HireInfor {
	// Handle potential conversion errors (e.g., time format)
	createdAt := hire.CreatedAt.Format(time.RFC3339Nano)
	deletedAt := ""
	updatedAt := ""

	if !hire.UpdatedAt.IsZero() {
		updatedAt = hire.UpdatedAt.Format(time.RFC3339Nano)
	}

	if !hire.DeletedAt.Time.IsZero() {
		deletedAt = hire.DeletedAt.Time.Format(time.RFC3339Nano)
	}
	return &proapi.HireInfor{
		Id:              hire.ID.String(),
		CreatedAt:       createdAt,
		CreatedBy:       hire.CreatedBy,
		UpdatedAt:       updatedAt,
		UpdatedBy:       hire.UpdatedBy,
		DeletedBy:       hire.DeletedBy,
		DeletedAt:       deletedAt,
		UserId:          hire.UserId,
		ProviderId:      hire.ProviderId.String(),
		ServiceId:       hire.ServiceId.String(),
		WorkTimeFrom:    hire.WorkTimeFrom,
		WorkTimeTo:      hire.WorkTimeTo,
		Status:          hire.Status,
		PaymentMethodId: hire.PaymentMethodId.String(),
		Issue:           hire.Issue,
		Service:         MapToService(hire.Service),
		Provider:        MapProviderToInfo(hire.Provider),
		Review:          MapToReview(hire.Review),
	}
}

func MapToListHireInfors(hires []*datasource.Hire) []*proapi.HireInfor {
	var listHire = make([]*proapi.HireInfor, 0)
	for _, hire := range hires {
		listHire = append(listHire, MapToHireInfor(hire))
	}
	return listHire
}
