package pro

import (
	"context"

	proapi "github.com/nguyentrunghieu15/be-beehome-prj/api/pro-api"
	"google.golang.org/protobuf/types/known/emptypb"
)

type HireAPIServer struct {
	proapi.UnimplementedHireServiceServer
}

func NewHireAPI() *HireAPIServer {
	return &HireAPIServer{}
}

func (s *HireAPIServer) CreateHire(
	ctx context.Context,
	req *proapi.CreateHireRequest,
) (*proapi.CreateHireResponse, error) {
	return nil, nil
}
func (s *HireAPIServer) DeleteHire(ctx context.Context, req *proapi.DeleteHireRequest) (*emptypb.Empty, error) {
	return nil, nil
}

func (s *HireAPIServer) FindAllHire(
	ctx context.Context,
	req *proapi.FindAllHireRequest,
) (*proapi.FindAllHireResponse, error) {
	return nil, nil
}

func (s *HireAPIServer) UpdateStatusHire(
	ctx context.Context,
	req *proapi.UpdateStatusHireRequest,
) (*proapi.UpdateStatusHireResponse, error) {
	return nil, nil
}
