package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	userapi "github.com/nguyentrunghieu15/be-beehome-prj/api/user-api"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/convert"
	"github.com/nguyentrunghieu15/be-beehome-prj/user-manager-service/internal/common"
	communication "github.com/nguyentrunghieu15/be-beehome-prj/user-manager-service/internal/comunitication"
	"github.com/nguyentrunghieu15/be-beehome-prj/user-manager-service/internal/mapper"
	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *UserService) BlockUser(ctx context.Context, req *userapi.BlockRequest) (*emptypb.Empty, error) {
	s.logger.Infor(common.StandardMsgInfor(ctx, "block user", "user id:"+req.Id))

	//validate date
	if err := s.validator.Validate(req); err != nil {
		s.logger.Error(common.StandardMsgError(ctx, "block user", err))
		return nil, status.Error(codes.InvalidArgument, "block user fail by invalid data")
	}

	// parse
	mapReq, err := convert.StructProtoToMap(req)
	if err != nil {
		s.logger.Error(common.StandardMsgError(ctx, "block user", err))
		return nil, status.Error(codes.Internal, "internal server")
	}

	// ban
	if _, err := s.bannedAccountRepo.CreateBannedAccount(mapReq); err != nil {
		s.logger.Error(common.StandardMsgError(ctx, "block user", err))
		return nil, status.Error(codes.Internal, "internal server")
	}

	// ban
	if _, err := s.userRepo.UpdateOneById(uuid.MustParse(req.Id), map[string]interface{}{"status": "banned"}); err != nil {
		s.logger.Error(common.StandardMsgError(ctx, "block user", err))
		return nil, status.Error(codes.Internal, "internal server")
	}

	// User successfully blocked, return empty response
	return &emptypb.Empty{}, nil
}

func (s *UserService) CreateUser(ctx context.Context, req *userapi.CreateUserRequest) (*userapi.UserInfor, error) {
	s.logger.Infor(common.StandardMsgInfor(ctx, "create user", "email:"+req.Email))

	// permision

	//validate date
	if err := s.validator.Validate(req); err != nil {
		s.logger.Error(common.StandardMsgError(ctx, "create user", err))
		return nil, status.Error(codes.InvalidArgument, "create user fail by invalid data")
	}

	// get user

	user, _ := s.userRepo.FindOneByEmail(req.Email)
	if user != nil {
		s.logger.Error(common.StandardMsgError(ctx, "create user", errors.New("exsit email:"+req.Email)))
		return nil, status.Error(codes.Internal, "internal server")
	}

	//parse to map
	mapUser, err := convert.StructProtoToMap(req)
	if err != nil {
		s.logger.Error(common.StandardMsgError(ctx, "create user", err))
		return nil, status.Error(codes.Internal, "internal server")
	}
	mapUser["status"] = "active"
	// Use userRepo to create a new user based on the request information
	newUser, err := s.userRepo.CreateUser(mapUser)
	if err != nil {
		s.logger.Error(common.StandardMsgError(ctx, "create user", err))
		return nil, status.Error(codes.Internal, "internal server")
	}

	tranferMsg, err := json.Marshal(map[string]interface{}{
		"type":    "create",
		"user_id": newUser.ID.String(),
		"role":    "user",
	})
	if err != nil {
		return nil, err
	}
	err = communication.UserResourceKafka.WriteMessages(
		context.Background(),
		kafka.Message{
			Value: tranferMsg,
		},
	)

	userInfo, err := mapper.ConvertUserToUserInfor(newUser)
	if err != nil {
		s.logger.Error(common.StandardMsgError(ctx, "create user", err))
		return nil, status.Error(codes.Internal, "internal server")
	}

	return userInfo, nil
}

// ListUsers retrieves a list of users based on request parameters
func (s *UserService) ListUsers(
	ctx context.Context,
	req *userapi.ListUsersRequest,
) (*userapi.ListUsersResponse, error) {
	s.logger.Infor(common.StandardMsgInfor(ctx, "list users", fmt.Sprintf("filter by: %+v", req)))

	// Validate user ID format (optional, can be done before)
	if err := s.validator.Validate(req); err != nil {
		s.logger.Error(common.StandardMsgError(ctx, "list users", errors.New("invalid data format")))
		return nil, status.Error(codes.InvalidArgument, "invalid data format")
	}

	// Use userRepo to retrieve a list of users based on the request filters
	users, err := s.userRepo.FindUsers(req)
	if err != nil {
		s.logger.Error(common.StandardMsgError(ctx, "list users", err))
		return nil, status.Error(codes.Internal, "internal server error")
	}

	// Convert users to UserInfor format
	userInfos, err := mapper.ConvertListUserToListUserInfor(users)
	if err != nil {
		s.logger.Error(common.StandardMsgError(ctx, "list users", err))
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &userapi.ListUsersResponse{Users: userInfos}, nil
}

// GetUser retrieves a user by its ID
func (s *UserService) GetUser(ctx context.Context, req *userapi.GetUserRequest) (*userapi.UserInfor, error) {
	s.logger.Infor(common.StandardMsgInfor(ctx, "get user", fmt.Sprintf("user id: %s", req.GetId())))

	// Validate user ID format (optional, can be done before)
	if err := s.validator.Validate(req); err != nil {
		s.logger.Error(common.StandardMsgError(ctx, "get user", errors.New("invalid user ID format")))
		return nil, status.Error(codes.InvalidArgument, "invalid user ID format")
	}

	// Use userRepo to retrieve a user by its ID
	user, err := s.userRepo.FindOneById(uuid.MustParse(req.GetId()))
	if err != nil {
		s.logger.Error(common.StandardMsgError(ctx, "get user", err))
		return nil, status.Error(codes.Internal, "not found user id")
	}

	// Convert user to UserInfor format
	userInfo, err := mapper.ConvertUserToUserInfor(user)
	if err != nil {
		s.logger.Error(common.StandardMsgError(ctx, "get user", err))
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return userInfo, nil
}

// UpdateUser updates an existing user
func (s *UserService) UpdateUser(ctx context.Context, req *userapi.UpdateUserRequest) (*userapi.UserInfor, error) {
	s.logger.Infor(common.StandardMsgInfor(ctx, "update user", "user id:"+req.Id))

	// validate data
	if err := s.validator.Validate(req); err != nil {
		s.logger.Error(common.StandardMsgError(ctx, "update user", err))
		return nil, status.Error(codes.InvalidArgument, "update user fail by invalid data")
	}

	// Use userRepo to retrieve a user by its ID
	_, err := s.userRepo.FindOneById(uuid.MustParse(req.GetId()))
	if err != nil {
		s.logger.Error(common.StandardMsgError(ctx, "update user", err))
		return nil, status.Error(codes.Internal, "fail to find your account")
	}

	// parse to map
	mapReq, err := convert.StructProtoToMap(req)
	if err != nil {
		s.logger.Error(common.StandardMsgError(ctx, "update user", err))
		return nil, status.Error(codes.Internal, "internal server")
	}

	// update user using userRepo
	updatedUser, err := s.userRepo.UpdateOneById(uuid.MustParse(req.Id), mapReq)
	if err != nil {
		s.logger.Error(common.StandardMsgError(ctx, "update user", err))
		return nil, status.Error(codes.Internal, "internal server")
	}

	userInfo, err := mapper.ConvertUserToUserInfor(updatedUser)
	if err != nil {
		s.logger.Error(common.StandardMsgError(ctx, "update user", err))
		return nil, status.Error(codes.Internal, "internal server")
	}

	tranferMsg, err := json.Marshal(map[string]interface{}{
		"type":        "delete",
		"user_id":     req.Id,
		"role":        "user",
		"provider_id": updatedUser.ProviderId,
	})
	if err != nil {
		return nil, err
	}
	communication.UserResourceKafka.WriteMessages(
		context.Background(),
		kafka.Message{
			Value: tranferMsg,
		},
	)

	return userInfo, nil
}

// DeleteUser deletes a user
func (s *UserService) DeleteUser(ctx context.Context, req *userapi.DeleteUserRequest) (*emptypb.Empty, error) {
	s.logger.Infor(common.StandardMsgInfor(ctx, "delete user", "user id:"+req.Id))

	// validate data
	if err := s.validator.Validate(req); err != nil {
		s.logger.Error(common.StandardMsgError(ctx, "delete user", err))
		return nil, status.Error(codes.InvalidArgument, "delete user fail by invalid data")
	}

	// ban
	if _, err := s.userRepo.UpdateOneById(uuid.MustParse(req.Id), map[string]interface{}{"status": "deactive"}); err != nil {
		s.logger.Error(common.StandardMsgError(ctx, "block user", err))
		return nil, status.Error(codes.Internal, "internal server")
	}

	// delete user using userRepo
	err := s.userRepo.DeleteOneById(uuid.MustParse(req.Id))
	if err != nil {
		s.logger.Error(common.StandardMsgError(ctx, "delete user", err))
		return nil, status.Error(codes.Internal, "internal server")
	}

	tranferMsg, err := json.Marshal(map[string]interface{}{
		"type":    "delete",
		"user_id": req.Id,
		"role":    "user",
	})
	if err != nil {
		return nil, err
	}
	communication.UserResourceKafka.WriteMessages(
		context.Background(),
		kafka.Message{
			Value: tranferMsg,
		},
	)

	// User successfully deleted, return empty response
	return &emptypb.Empty{}, nil
}
