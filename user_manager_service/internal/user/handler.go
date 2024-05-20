package user

import (
	"context"

	userapi "github.com/nguyentrunghieu15/be-beehome-prj/api/user-api"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/convert"
	"github.com/nguyentrunghieu15/be-beehome-prj/user_manager_service/internal/common"
	"github.com/nguyentrunghieu15/be-beehome-prj/user_manager_service/internal/mapper"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *UserService) BlockUser(ctx context.Context, req *userapi.BlockRequest) (*emptypb.Empty, error) {

	// User successfully blocked, return empty response
	return &emptypb.Empty{}, nil
}

func (s *UserService) CreateUser(ctx context.Context, req *userapi.CreateUserRequest) (*userapi.UserInfor, error) {
	// permision

	//validate date
	if err := s.validator.Validate(req); err != nil {
		s.logger.Error(common.StandardMsgError(ctx, "create user", err))
		return nil, status.Error(codes.InvalidArgument, "create user fail by invalid data")
	}

	//parse to map
	mapUser, err := convert.StructProtoToMap(req)
	if err != nil {
		s.logger.Error(common.StandardMsgError(ctx, "create user", err))
		return nil, status.Error(codes.Internal, "internal server")
	}
	// Use userRepo to create a new user based on the request information
	newUser, err := s.userRepo.CreateUser(mapUser)
	if err != nil {
		s.logger.Error(common.StandardMsgError(ctx, "create user", err))
		return nil, status.Error(codes.Internal, "internal server")
	}

	userInfo, err := mapper.ConvertUserToUserInfor(*newUser)
	if err != nil {
		s.logger.Error(common.StandardMsgError(ctx, "create user", err))
		return nil, status.Error(codes.Internal, "internal server")
	}

	return userInfo, nil
}
