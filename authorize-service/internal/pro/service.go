package pro

import (
	"context"

	proapi "github.com/nguyentrunghieu15/be-beehome-prj/api/pro-api"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ServiceAPIServer struct {
	proapi.UnimplementedServiceManagerServiceServer
}

func NewServiceAPI() *ServiceAPIServer {
	return &ServiceAPIServer{}
}

func (s *ServiceAPIServer) CreateService(
	ctx context.Context,
	req *proapi.CreateServiceRequest,
) (*proapi.Service, error) {
	return nil, nil
}

func (s *ServiceAPIServer) DeleteService(
	ctx context.Context,
	req *proapi.DeleteServiceRequest,
) (*emptypb.Empty, error) {
	return nil, nil
}

func (s *ServiceAPIServer) FulltextSearchServices(
	ctx context.Context,
	req *proapi.FulltextSearchServicesRequest,
) (*proapi.ListServicesResponse, error) {
	return nil, nil
}
func (s *ServiceAPIServer) GetService(ctx context.Context, req *proapi.GetServiceRequest) (*proapi.Service, error) {
	return nil, nil
}

func (s *ServiceAPIServer) ListServices(
	ctx context.Context,
	req *proapi.ListServicesRequest,
) (*proapi.ListServicesResponse, error) {
	return nil, nil
}

func (s *ServiceAPIServer) UpdateService(
	ctx context.Context,
	req *proapi.UpdateServiceRequest,
) (*proapi.Service, error) {
	return nil, nil
}
