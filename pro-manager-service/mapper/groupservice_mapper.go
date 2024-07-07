package mapper

import (
	"time"

	proapi "github.com/nguyentrunghieu15/be-beehome-prj/api/pro-api"
	"github.com/nguyentrunghieu15/be-beehome-prj/pro-manager-service/internal/datasource"
)

func MapToGroupService(ds *datasource.GroupService) *proapi.GroupService {
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
	return &proapi.GroupService{
		Id:        id,
		CreatedAt: createdAt,
		CreatedBy: ds.CreatedBy,
		UpdatedAt: updatedAt,
		UpdatedBy: ds.UpdatedBy,
		DeletedBy: ds.DeletedBy,
		DeletedAt: deletedAt,
		Name:      ds.Name,
		Detail:    ds.Detail,
		Services:  MapToServices(ds.Services),
	}
}

func MapToGroupServices(ds []*datasource.GroupService) []*proapi.GroupService {
	proapiGroupServices := make([]*proapi.GroupService, 0, len(ds))
	for _, gs := range ds {
		proapiGroupServices = append(proapiGroupServices, MapToGroupService(gs))
	}
	return proapiGroupServices
}
