package profiles

import (
	"context"
	"encoding/json"
	"errors"

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
	"gorm.io/gorm"
)

func (ps *ProfileService) AddCard(ctx context.Context, req *userapi.AddCardRequest) (*emptypb.Empty, error) {
	ps.logger.Infor(common.StandardMsgInfor(ctx, "add card", "user id:"+ctx.Value("user_id").(string)))
	// Validate the request using the validator
	if err := ps.validator.Validate(req); err != nil {
		ps.logger.Error(common.StandardMsgError(ctx, "add card", err))
		return nil, status.Error(codes.InvalidArgument, "invalid data")
	}

	// Extract card information from the request
	cardInfo := req.Card // extract card details from req
	mapCard, err := convert.StructProtoToMap(cardInfo)
	if err != nil {
		ps.logger.Error(common.StandardMsgError(ctx, "add card", err))
		return nil, status.Error(codes.Internal, "internal server")
	}

	mapCard["user_id"] = ctx.Value("user_id")

	// Use the cardRepo to add the card to the user's profile
	if _, err := ps.cardRepo.CreateCard(mapCard); err != nil {
		ps.logger.Error(common.StandardMsgError(ctx, "add card", err))
		return nil, status.Error(codes.Internal, "internal server")
	}

	// Return an empty response if successful
	return &emptypb.Empty{}, nil
}

func (ps *ProfileService) ChangeEmail(
	ctx context.Context,
	req *userapi.ChangeEmailRequest,
) (*userapi.UserInfor, error) {
	ps.logger.Infor(common.StandardMsgInfor(ctx, "change email", "user id:"+ctx.Value("user_id").(string)))
	// Validate the request using the validator
	if err := ps.validator.Validate(req); err != nil {
		ps.logger.Error(common.StandardMsgError(ctx, "change email", err))
		return nil, status.Error(codes.InvalidArgument, "invalid data")
	}

	// Extract new email from the request
	newEmail := req.Email
	mapData := map[string]interface{}{
		"email": newEmail,
	}

	// Check if the new email already exists in the database
	_, err := ps.userRepo.FindOneByEmail(newEmail)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		ps.logger.Error(common.StandardMsgError(ctx, "change email", err))
		return nil, status.Error(codes.Internal, "internal server")
	}

	// Use the userRepo to update the user's email
	user, err := ps.userRepo.UpdateOneById(uuid.MustParse(ctx.Value("user_id").(string)), mapData)
	if err != nil {
		ps.logger.Error(common.StandardMsgError(ctx, "change email", err))
		return nil, status.Error(codes.Internal, "internal server")
	}

	userInfor, err := mapper.ConvertUserToUserInfor(user)
	if err != nil {
		ps.logger.Error(common.StandardMsgError(ctx, "change email", err))
		return nil, status.Error(codes.Internal, "internal server")
	}
	// Return the updated user information
	return userInfor, nil
}

func (ps *ProfileService) ChangeName(ctx context.Context, req *userapi.ChangeNameRequest) (*userapi.UserInfor, error) {
	ps.logger.Infor(common.StandardMsgInfor(ctx, "change name", "user id:"+ctx.Value("user_id").(string)))
	// Validate the request using the validator
	if err := ps.validator.Validate(req); err != nil {
		ps.logger.Error(common.StandardMsgError(ctx, "change name", err))
		return nil, err
	}

	// Extract new name from the request
	newName := req.GetName()
	mapData := map[string]interface{}{
		"first_name": newName,
	}

	// Update user's name
	user, err := ps.userRepo.UpdateOneById(uuid.MustParse(ctx.Value("user_id").(string)), mapData)
	if err != nil {
		// Consider adding more specific error handling based on your userRepo implementation
		ps.logger.Error(common.StandardMsgError(ctx, "change name", err))
		return nil, status.Error(codes.Internal, "internal server error")
	}

	userInfo, err := mapper.ConvertUserToUserInfor(user)
	if err != nil {
		ps.logger.Error(common.StandardMsgError(ctx, "change name", err))
		return nil, status.Error(codes.Internal, "internal server")
	}
	// Return the updated user information
	return userInfo, nil
}

func (ps *ProfileService) GetProfile(ctx context.Context, req *emptypb.Empty) (*userapi.UserInfor, error) {
	ps.logger.Infor(common.StandardMsgInfor(ctx, "get profile", "user id:"+ctx.Value("user_id").(string)))
	// Retrieve user information from the userRepo
	user, err := ps.userRepo.FindOneProfileById(uuid.MustParse(ctx.Value("user_id").(string)))
	if err != nil {
		// Consider adding more specific error handling based on your userRepo implementation
		ps.logger.Error(common.StandardMsgError(ctx, "get profile", err))
		return nil, status.Error(codes.Internal, "internal server error")
	}

	// Optionally, retrieve additional information (cards?) from cardRepo
	userInfo, err := mapper.ConvertUserToUserInfor(user)
	if err != nil {
		ps.logger.Error(common.StandardMsgError(ctx, "change name", err))
		return nil, status.Error(codes.Internal, "internal server")
	}
	// Return the user information
	return userInfo, nil
}

func (ps *ProfileService) DeactiveAccount(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	userId, err := uuid.Parse(ctx.Value("user_id").(string))
	if err != nil {
		ps.logger.Error(common.StandardMsgError(ctx, "deactive account", err))
		return nil, status.Error(codes.Internal, "internal server")
	}

	_, err = ps.userRepo.UpdateOneById(userId, map[string]interface{}{"status": "deactive"})
	if err != nil {
		ps.logger.Error(common.StandardMsgError(ctx, "deactive account", err))
		return nil, status.Error(codes.Internal, "internal server")
	}
	err = ps.userRepo.DeleteOneById(userId)
	if err != nil {
		ps.logger.Error(common.StandardMsgError(ctx, "deactive account", err))
		return nil, status.Error(codes.Internal, "internal server")
	}
	tranferMsg, err := json.Marshal(map[string]interface{}{
		"type":    "delete",
		"user_id": userId.String(),
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
	return &emptypb.Empty{}, nil
}
