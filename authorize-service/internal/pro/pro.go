package pro

import (
	"context"

	proapi "github.com/nguyentrunghieu15/be-beehome-prj/api/pro-api"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ProAPIServer struct {
	proapi.UnimplementedProServiceServer
}

func NewProAPI() *ProAPIServer {
	return &ProAPIServer{}
}

func (s *ProAPIServer) AddPaymentMethodPro(
	ctx context.Context,
	req *proapi.AddPaymentMethodProRequest,
) (*emptypb.Empty, error) {
	return nil, nil
}
func (s *ProAPIServer) AddServicePro(ctx context.Context, req *proapi.AddServiceProRequest) (*emptypb.Empty, error) {
	return nil, nil
}

func (s *ProAPIServer) AddSocialMediaPro(
	ctx context.Context,
	req *proapi.AddSocialMediaProRequest,
) (*emptypb.Empty, error) {
	return nil, nil
}
func (s *ProAPIServer) DeleteProById(ctx context.Context, req *proapi.DeleteProByIdRequest) (*emptypb.Empty, error) {
	return nil, nil
}

func (s *ProAPIServer) DeleteServicePro(
	ctx context.Context,
	req *proapi.DeleteServiceProRequest,
) (*emptypb.Empty, error) {
	return nil, nil
}

func (s *ProAPIServer) DeleteSocialMediaPro(
	ctx context.Context,
	req *proapi.DeleteSocialMediaProRequest,
) (*emptypb.Empty, error) {
	return nil, nil
}

func (s *ProAPIServer) FindProById(
	ctx context.Context,
	req *proapi.FindProByIdRequest,
) (*proapi.FindProByIdResponse, error) {
	return nil, nil
}
func (s *ProAPIServer) FindPros(ctx context.Context, req *proapi.FindProsRequest) (*proapi.FindProsResponse, error) {
	return nil, nil
}

func (s *ProAPIServer) GetAllReviewsOfProvider(
	ctx context.Context,
	req *proapi.GetAllReviewOfProviderRequest,
) (*proapi.GetAllReviewOfProviderResponse, error) {
	return nil, nil
}

func (s *ProAPIServer) GetAllServiceOfProvider(
	ctx context.Context,
	req *proapi.GetAllServiceOfProviderRequest,
) (*proapi.GetAllServiceOfProviderResponse, error) {
	return nil, nil
}

func (s *ProAPIServer) GetProviderProfile(
	ctx context.Context,
	req *emptypb.Empty,
) (*proapi.ProviderProfileResponse, error) {
	return nil, nil
}

func (s *ProAPIServer) JoinAsProvider(
	ctx context.Context,
	req *proapi.JoinAsProviderRequest,
) (*proapi.JoinAsProviderResponse, error) {
	return nil, nil
}
func (s *ProAPIServer) ReplyReviewPro(ctx context.Context, req *proapi.ReplyReviewProRequest) (*emptypb.Empty, error) {
	return nil, nil
}
func (s *ProAPIServer) ReviewPro(ctx context.Context, req *proapi.ReviewProRequest) (*emptypb.Empty, error) {
	return nil, nil
}
func (s *ProAPIServer) SignUpPro(ctx context.Context, req *proapi.SignUpProRequest) (*proapi.ProviderInfo, error) {
	return nil, nil
}
func (s *ProAPIServer) UpdatePro(ctx context.Context, req *proapi.UpdateProRequest) (*proapi.ProviderInfo, error) {
	return nil, nil
}

func (s *ProAPIServer) UpdateSocialMediaPro(
	ctx context.Context,
	req *proapi.UpdateSocialMediaProRequest,
) (*emptypb.Empty, error) {
	return nil, nil
}
