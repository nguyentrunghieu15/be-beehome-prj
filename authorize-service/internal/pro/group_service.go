package pro

import (
	"context"

	proapi "github.com/nguyentrunghieu15/be-beehome-prj/api/pro-api"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GroupServiceAPIServer struct {
	proapi.UnimplementedGroupServiceManagerServer
}

func NewGroupServiceAPI() *GroupServiceAPIServer {
	return &GroupServiceAPIServer{}
}

func (s *GroupServiceAPIServer) CreateGroupService(
	ctx context.Context,
	req *proapi.CreateGroupServiceRequest,
) (*proapi.GroupService, error) {
	return nil, nil
}

func (s *GroupServiceAPIServer) DeleteGroupService(
	ctx context.Context,
	req *proapi.DeleteGroupServiceRequest,
) (*emptypb.Empty, error) {
	return nil, nil
}

func (s *GroupServiceAPIServer) FulltextSearchGroupServices(
	ctx context.Context,
	req *proapi.FulltextSearchGroupServicesRequest,
) (*proapi.ListGroupServicesResponse, error) {
	return nil, nil
}

func (s *GroupServiceAPIServer) GetGroupService(
	ctx context.Context,
	req *proapi.GetGroupServiceRequest,
) (*proapi.GroupService, error) {
	return nil, nil
}

func (s *GroupServiceAPIServer) ListGroupServices(
	ctx context.Context,
	req *proapi.ListGroupServicesRequest,
) (*proapi.ListGroupServicesResponse, error) {
	return nil, nil
}

func (s *GroupServiceAPIServer) UpdateGroupService(
	ctx context.Context,
	req *proapi.UpdateGroupServiceRequest,
) (*proapi.GroupService, error) {
	return nil, nil
}
