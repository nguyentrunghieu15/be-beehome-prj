package mapper

import (
	"time"

	proapi "github.com/nguyentrunghieu15/be-beehome-prj/api/pro-api"
	"github.com/nguyentrunghieu15/be-beehome-prj/pro-manager-service/internal/datasource"
)

func MapToHire(hire *datasource.Hire) *proapi.Hire {
	return &proapi.Hire{
		Id:              hire.ID.String(),
		CreatedAt:       hire.CreatedAt.Format(time.RFC3339Nano),
		CreatedBy:       hire.CreatedBy,
		UpdatedAt:       hire.UpdatedAt.Format(time.RFC3339Nano),
		UpdatedBy:       hire.UpdatedBy,
		DeletedBy:       hire.DeletedBy,
		DeletedAt:       hire.DeletedAt.Time.Format(time.RFC3339Nano),
		UserId:          hire.UserId,
		ProviderId:      hire.ProviderId.String(),
		ServiceId:       hire.ServiceId.String(),
		WorkTimeFrom:    hire.WorkTimeFrom,
		WorkTimeTo:      hire.WorkTimeTo,
		Status:          hire.Status,
		PaymentMethodId: hire.PaymentMethodId.String(),
	}
}

func MapToListHire(hires []*datasource.Hire) []*proapi.Hire {
	var listHire = make([]*proapi.Hire, 0)
	for _, hire := range hires {
		listHire = append(listHire, MapToHire(hire))
	}
	return listHire
}
