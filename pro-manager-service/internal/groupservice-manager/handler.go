package groupservicemanager

import (
	"context"

	"github.com/google/uuid"
	proapi "github.com/nguyentrunghieu15/be-beehome-prj/api/pro-api"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/convert"
	"github.com/nguyentrunghieu15/be-beehome-prj/pro-manager-service/mapper"
	"google.golang.org/protobuf/types/known/emptypb"
)

// ListServices implements the ListServices RPC method.
func (gs *GroupServiceManagerServer) FulltextSearchGroupServices(
	ctx context.Context,
	req *proapi.FulltextSearchGroupServicesRequest,
) (*proapi.ListGroupServicesResponse, error) {
	// Validate GetServiceRequest (e.g., check if ID is empty)
	if err := gs.validator.Validate(req); err != nil {
		return nil, err
	}

	groups, err := gs.groupServiceRepo.FulltextSearchGroupServiceByName(req.Name)
	if err != nil {
		return nil, err
	}

	return &proapi.ListGroupServicesResponse{GroupServices: mapper.MapToGroupServices(groups)}, nil
}

// GetService implements the GetService RPC method.
func (gs *GroupServiceManagerServer) GetGroupService(
	ctx context.Context,
	req *proapi.GetGroupServiceRequest,
) (*proapi.GroupService, error) {
	// Validate GetServiceRequest (e.g., check if ID is empty)
	if err := gs.validator.Validate(req); err != nil {
		return nil, err
	}

	gpServiceId := uuid.MustParse(req.Id)

	gpService, err := gs.groupServiceRepo.FindOneById(gpServiceId)
	if err != nil {
		gs.logger.Errorf("failed to get service: %v", err)
		return nil, err
	}

	return mapper.MapToGroupService(gpService), nil
}

// CreateService implements the CreateService RPC method.
func (gs *GroupServiceManagerServer) CreateGroupService(
	ctx context.Context,
	req *proapi.CreateGroupServiceRequest,
) (*proapi.GroupService, error) {
	// Validate CreateServiceRequest (e.g., check if name and group_id are empty)
	if err := gs.validator.Validate(req); err != nil {
		return nil, err
	}

	mapData, err := convert.StructProtoToMap(req)
	if err != nil {
		gs.logger.Errorf("failed to convert service: %v", err)
		return nil, err
	}

	gpService, err := gs.groupServiceRepo.CreateGroupService(mapData)
	if err != nil {
		gs.logger.Errorf("failed to create service: %v", err)
		return nil, err
	}

	return mapper.MapToGroupService(gpService), nil
}

// UpdateService implements the UpdateService RPC method.
func (gs *GroupServiceManagerServer) UpdateGroupService(
	ctx context.Context,
	req *proapi.UpdateGroupServiceRequest,
) (*proapi.GroupService, error) {
	// Validate UpdateServiceRequest (e.g., check if ID and potentially other fields are empty)
	if err := gs.validator.Validate(req); err != nil {
		return nil, err
	}
	mapData, err := convert.StructProtoToMap(req)
	if err != nil {
		gs.logger.Errorf("failed to convert service: %v", err)
		return nil, err
	}
	delete(mapData, "id")

	gpServiceId := uuid.MustParse(req.Id)

	gpService, err := gs.groupServiceRepo.UpdateOneById(gpServiceId, mapData)
	if err != nil {
		gs.logger.Errorf("failed to update service: %v", err)
		return nil, err
	}

	return mapper.MapToGroupService(gpService), nil
}

// DeleteService implements the DeleteService RPC method.
func (gs *GroupServiceManagerServer) DeleteService(
	ctx context.Context,
	req *proapi.DeleteServiceRequest,
) (*emptypb.Empty, error) {
	// Validate DeleteServiceRequest (e.g., check if ID is empty)
	if err := gs.validator.Validate(req); err != nil {
		return nil, err
	}

	gpServiceId := uuid.MustParse(req.Id)
	err := gs.groupServiceRepo.DeleteOneById(gpServiceId)
	if err != nil {
		gs.logger.Errorf("failed to delete service: %v", err)
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (gs *GroupServiceManagerServer) ListGroupServices(
	ctx context.Context,
	req *proapi.ListGroupServicesRequest,
) (*proapi.ListGroupServicesResponse, error) {
	// Validate DeleteServiceRequest (e.g., check if ID is empty)
	if err := gs.validator.Validate(req); err != nil {
		return nil, err
	}

	gpServices, err := gs.groupServiceRepo.FindGroupServices(req)
	if err != nil {
		return nil, err
	}

	return &proapi.ListGroupServicesResponse{GroupServices: mapper.MapToGroupServices(gpServices)}, nil
}
