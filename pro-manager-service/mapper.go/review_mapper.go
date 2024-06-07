package mapper

import (
	"time"

	proapi "github.com/nguyentrunghieu15/be-beehome-prj/api/pro-api"
	"github.com/nguyentrunghieu15/be-beehome-prj/pro-manager-service/internal/datasource"
)

func MapToReview(ds *datasource.Review) *proapi.Review {
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
	userId := ds.UserId.String()
	providerId := ds.ProviderId.String()
	serviceId := ds.ServiceId.String() // Assuming ServiceId can be represented by int32

	// Map fields with matching names
	return &proapi.Review{
		Id:         id,
		CreatedAt:  createdAt,
		CreatedBy:  ds.CreatedBy,
		UpdatedAt:  updatedAt,
		UpdatedBy:  ds.UpdatedBy,
		DeletedBy:  ds.DeletedBy,
		DeletedAt:  deletedAt,
		UserId:     userId,
		ProviderId: providerId,
		Rating:     ds.Rating,
		Comment:    ds.Comment,
		Reply:      ds.Reply,
		ServiceId:  serviceId,
	}
}

func MapToReviews(ds []*datasource.Review) []*proapi.Review {
	proapiReviews := make([]*proapi.Review, 0, len(ds))
	for _, p := range ds {
		proapiReviews = append(proapiReviews, MapToReview(p))
	}
	return proapiReviews
}
