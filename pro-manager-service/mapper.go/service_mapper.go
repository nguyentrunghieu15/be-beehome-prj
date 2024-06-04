package mapper

import (
	"time"

	proapi "github.com/nguyentrunghieu15/be-beehome-prj/api/pro-api"
	"github.com/nguyentrunghieu15/be-beehome-prj/pro-manager-service/internal/datasource"
)

func MapToService(s *datasource.Service) *proapi.Service {
	// Handle potential conversion errors (e.g., time format)
	createdAt := s.CreatedAt.Format(time.RFC3339Nano)
	deletedAt := ""
	updatedAt := ""

	if !s.UpdatedAt.IsZero() {
		updatedAt = s.UpdatedAt.Format(time.RFC3339Nano)
	}

	if !s.DeletedAt.Time.IsZero() {
		deletedAt = s.DeletedAt.Time.Format(time.RFC3339Nano)
	}
	return &proapi.Service{
		Id:        s.ID.String(),
		CreatedAt: createdAt, // Assuming desired format for CreatedAt
		CreatedBy: s.CreatedBy,
		UpdatedAt: updatedAt, // Assuming desired format for UpdatedAt
		UpdatedBy: s.UpdatedBy,
		DeletedBy: s.DeletedBy,
		DeletedAt: deletedAt, // Assuming desired format for DeletedAt
		Name:      s.Name,
		Detail:    s.Detail,
		GroupId:   s.GroupServiceId.String(),
	}
}

func MapToServices(services []*datasource.Service) []*proapi.Service {
	proapiServices := make([]*proapi.Service, 0, len(services))
	for _, service := range services {
		proapiServices = append(proapiServices, MapToService(service))
	}
	return proapiServices
}
