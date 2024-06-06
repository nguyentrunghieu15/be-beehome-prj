package servicemanager

import (
	"context"

	"github.com/google/uuid"
	proapi "github.com/nguyentrunghieu15/be-beehome-prj/api/pro-api"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/convert"
	"github.com/nguyentrunghieu15/be-beehome-prj/pro-manager-service/mapper.go"
	"google.golang.org/protobuf/types/known/emptypb"
)

// ListServices implements the ListServices RPC method.
func (s *ServiceManagerServer) FulltextSearchServices(
	ctx context.Context,
	req *proapi.FulltextSearchServicesRequest,
) (*proapi.ListServicesResponse, error) {
	// Validate GetServiceRequest (e.g., check if ID is empty)
	if err := s.validator.Validate(req); err != nil {
		return nil, err
	}

	groups, err := s.groupServiceRepo.FulltextSearchGroupServiceByName(req.Name)
	if err != nil {
		return nil, err
	}

	groupIds := make([]string, 0)
	for _, g := range groups {
		groupIds = append(groupIds, g.ID.String())
	}

	services, err := s.serviceRepo.FulltextSearchServiceByNameOrInGroup(req.Name, groupIds...)
	if err != nil {
		return nil, err
	}

	return &proapi.ListServicesResponse{Services: mapper.MapToServices(services)}, nil
}

// GetService implements the GetService RPC method.
func (s *ServiceManagerServer) GetService(ctx context.Context, req *proapi.GetServiceRequest) (*proapi.Service, error) {
	// Validate GetServiceRequest (e.g., check if ID is empty)
	if err := s.validator.Validate(req); err != nil {
		return nil, err
	}

	serviceId := uuid.MustParse(req.Id)

	service, err := s.serviceRepo.FindOneById(serviceId)
	if err != nil {
		s.logger.Errorf("failed to get service: %v", err)
		return nil, err
	}

	return mapper.MapToService(service), nil
}

// CreateService implements the CreateService RPC method.
func (s *ServiceManagerServer) CreateService(
	ctx context.Context,
	req *proapi.CreateServiceRequest,
) (*proapi.Service, error) {
	// Validate CreateServiceRequest (e.g., check if name and group_id are empty)
	if err := s.validator.Validate(req); err != nil {
		return nil, err
	}

	mapData, err := convert.StructProtoToMap(req)
	if err != nil {
		s.logger.Errorf("failed to convert service: %v", err)
		return nil, err
	}

	service, err := s.serviceRepo.CreateService(mapData)
	if err != nil {
		s.logger.Errorf("failed to create service: %v", err)
		return nil, err
	}

	return mapper.MapToService(service), nil
}

// UpdateService implements the UpdateService RPC method.
func (s *ServiceManagerServer) UpdateService(
	ctx context.Context,
	req *proapi.UpdateServiceRequest,
) (*proapi.Service, error) {
	// Validate UpdateServiceRequest (e.g., check if ID and potentially other fields are empty)
	if err := s.validator.Validate(req); err != nil {
		return nil, err
	}
	mapData, err := convert.StructProtoToMap(req)
	if err != nil {
		s.logger.Errorf("failed to convert service: %v", err)
		return nil, err
	}
	delete(mapData, "id")

	serviceId := uuid.MustParse(req.Id)

	service, err := s.serviceRepo.UpdateOneById(serviceId, mapData)
	if err != nil {
		s.logger.Errorf("failed to update service: %v", err)
		return nil, err
	}

	return mapper.MapToService(service), nil
}

// DeleteService implements the DeleteService RPC method.
func (s *ServiceManagerServer) DeleteService(
	ctx context.Context,
	req *proapi.DeleteServiceRequest,
) (*emptypb.Empty, error) {
	// Validate DeleteServiceRequest (e.g., check if ID is empty)
	if err := s.validator.Validate(req); err != nil {
		return nil, err
	}

	serviceId := uuid.MustParse(req.Id)
	err := s.serviceRepo.DeleteOneById(serviceId)
	if err != nil {
		s.logger.Errorf("failed to delete service: %v", err)
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *ServiceManagerServer) ListServices(ctx context.Context, req *proapi.ListServicesRequest) (*proapi.ListServicesResponse, error) {
	// Validate DeleteServiceRequest (e.g., check if ID is empty)
	if err := s.validator.Validate(req); err != nil {
		return nil, err
	}

	services, err := s.serviceRepo.FindServices(req)
	if err != nil {
		return nil, err
	}
	return &proapi.ListServicesResponse{Services: mapper.MapToServices(services)}, nil
}
