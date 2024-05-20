package auth

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	pb "github.com/nguyentrunghieu15/be-beehome-prj/api/auth-api"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/convert"
	argon "github.com/nguyentrunghieu15/be-beehome-prj/internal/crypto/argon2"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/mail"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/random"
	"github.com/nguyentrunghieu15/be-beehome-prj/pkg/jwt"
	"github.com/nguyentrunghieu15/be-beehome-prj/user_manager_service/internal/common"
	"github.com/nguyentrunghieu15/be-beehome-prj/user_manager_service/internal/datasource"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	prefixAccsessTokenSessionStorage       = "login:session:access"
	prefixRefreshTokenSessionStorage       = "login:session:refresh"
	prefixForgotPasswordSessionStorage     = "forgotpassword:token"
	expireTimeForgotPasswordSessionStorage = 2 * time.Minute
)

// Login implements the Login RPC method
func (s *AuthService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	s.logger.Infor(
		common.StandardMsgInfor(ctx, "login", "email:"+req.Email),
	)
	// Implement validate infor
	if err := s.validator.Validate(req); err != nil {
		s.logger.Error(common.StandardMsgError(ctx, "login", err))
		return nil, status.Error(codes.InvalidArgument, "login fail by invalid data")
	}

	// This is a placeholder implementation
	user, err := s.userRepo.FindOneByEmail(req.Email)
	if err != nil {
		s.logger.Error(common.StandardMsgError(ctx, "login", err))
		return nil, status.Error(codes.Unauthenticated, "fail validate credentials")
	}

	if ok, err := argon.Compare(user.Password, req.Password); !ok {
		if err != nil {
			s.logger.Error(common.StandardMsgError(ctx, "login", err))
		} else {
			s.logger.Error(common.StandardMsgError(ctx, "login", errors.New("wrong password")))
		}
		return nil, status.Error(codes.Unauthenticated, "fail validate credentials")
	}

	var errs []error
	accessToken, err := s.jwtGenerator.GenerateToken(user.ID.String(),
		jwt.DefaultAccessTokenConfigure)
	if err != nil {
		errs = append(errs, err)
	}

	refresToken, err := s.jwtGenerator.GenerateToken(user.ID.String(),
		jwt.DefaultRefreshTokenConfigure)
	if err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		s.logger.Error(common.StandardMsgError(ctx, "login", err))
		return nil, status.Error(codes.Internal, "internal server")
	}

	if err := s.sessionStorage.SaveSession(
		fmt.Sprintf("%v:%v", prefixAccsessTokenSessionStorage, user.ID),
		accessToken,
		datasource.SessionKeyConfig{
			ExpireTime: jwt.DefaultAccessTokenConfigure.ExpiresTime,
		},
	); err != nil {
		s.logger.Error(common.StandardMsgError(ctx, "login", err))
		return nil, status.Error(codes.Internal, "internal server")
	}

	if err := s.sessionStorage.SaveSession(
		fmt.Sprintf("%v:%v", prefixRefreshTokenSessionStorage, user.ID),
		refresToken,
		datasource.SessionKeyConfig{
			ExpireTime: jwt.DefaultRefreshTokenConfigure.ExpiresTime,
		},
	); err != nil {
		s.logger.Error(common.StandardMsgError(ctx, "login", err))
		return nil, status.Error(codes.Internal, "internal server")
	}

	return &pb.LoginResponse{
		AccessToken:  accessToken,
		ExpireTime:   jwt.DefaultAccessTokenConfigure.ExpiresTime.Milliseconds(),
		RefreshToken: refresToken,
		TokenType:    "Bearer",
	}, nil
}

// RefreshToken implements the RefreshToken RPC method
func (s *AuthService) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	s.logger.Infor(
		common.StandardMsgInfor(ctx, "refresh token", ""),
	)
	// Implement your refresh token logic here (e.g., validate refresh token, generate new access token)
	// This is a placeholder implementation, replace with your actual logic
	if err := s.validator.Validate(req); err != nil {
		s.logger.Error(common.StandardMsgError(ctx, "refresh token", err))
		return nil, status.Error(codes.InvalidArgument, "refresh token fail by invalid data")
	}

	data, err := s.jwtGenerator.ParseToken(req.RefreshToken)
	if err != nil {
		s.logger.Error(common.StandardMsgError(ctx, "refresh token", err))
		return nil, status.Error(codes.Internal, "refresh token fail get data")
	}

	userId, ok := data.(string)
	if !ok {
		s.logger.Error(common.StandardMsgError(ctx, "refresh token", errors.New("cant get claims")))
		return nil, status.Error(codes.Internal, "refresh token fail get data")
	}

	//check exist refresh token
	if value, err := s.sessionStorage.GetSession(
		fmt.Sprintf("%v:%v", prefixRefreshTokenSessionStorage, userId),
	); value != req.RefreshToken || err != nil {
		s.logger.Error(common.StandardMsgError(ctx, "refresh token", errors.New("token not exist")))
		return nil, status.Error(codes.Internal, "refresh token fail get data")
	}

	newAccessToken, err := s.jwtGenerator.GenerateToken(userId, jwt.DefaultAccessTokenConfigure)
	if err != nil {
		s.logger.Error(common.StandardMsgError(ctx, "refresh token", err))
		return nil, status.Error(codes.Internal, "internal server")
	}

	if s.sessionStorage.SaveSession(
		fmt.Sprintf("%v:%v", prefixAccsessTokenSessionStorage, userId),
		newAccessToken,
		datasource.SessionKeyConfig{
			ExpireTime: jwt.DefaultAccessTokenConfigure.ExpiresTime,
		},
	) != nil {
		s.logger.Error(common.StandardMsgError(ctx, "refresh token", errors.New("cant save session token")))
		return nil, status.Error(codes.Internal, "internal server")
	}

	return &pb.RefreshTokenResponse{
		AccessToken: newAccessToken,
	}, nil
}

// ForgotPassword implements the ForgotPassword RPC method
func (s *AuthService) ForgotPassword(ctx context.Context, req *pb.ForgotPasswordRequest) (*emptypb.Empty, error) {
	s.logger.Infor(
		common.StandardMsgInfor(ctx, "forgot password", "email:"+req.Email),
	)
	// Implement your forgot password logic here (e.g., send reset password email)
	// This is a placeholder implementation, replace with your actual logic
	if err := s.validator.Validate(req); err != nil {
		s.logger.Error(common.StandardMsgError(ctx, "forgot password", err))
		return nil, status.Error(codes.InvalidArgument, "forgot password fail by invalid data")
	}

	_, err := s.userRepo.FindOneByEmail(req.Email)
	if err != nil {
		s.logger.Error(common.StandardMsgError(ctx, "forgot password", errors.New("not found email")))
		return &emptypb.Empty{}, nil
	}

	randToken, err := random.GenerateRandomBytes(40)
	if err != nil {
		s.logger.Error(common.StandardMsgError(ctx, "forgot password", err))
		return nil, status.Error(codes.Internal, "internal server")
	}

	tokenResetPassword := base64.RawStdEncoding.EncodeToString(randToken)

	if err := s.sessionStorage.SaveSession(
		fmt.Sprintf("%v:%v", prefixForgotPasswordSessionStorage, tokenResetPassword),
		req.Email,
		datasource.SessionKeyConfig{
			ExpireTime: expireTimeForgotPasswordSessionStorage,
		},
	); err != nil {
		s.logger.Error(common.StandardMsgError(ctx, "forgot password", err))
		return nil, status.Error(codes.Internal, "internal server")
	}

	if err := s.mailService.SendMail(mail.AuthStmp{}, mail.Letter{}); err != nil {
		s.logger.Error(common.StandardMsgError(ctx, "forgot password", err))
		return nil, status.Error(codes.Internal, "can't send mail")
	}
	return &emptypb.Empty{}, nil
}

// ResetPassword implements the ResetPassword RPC method
func (s *AuthService) ResetPassword(ctx context.Context, req *pb.ResetPasswordRequest) (*emptypb.Empty, error) {
	s.logger.Infor(
		common.StandardMsgInfor(ctx, "reset password", ""),
	)
	// Implement your reset password logic here (e.g., validate reset token, update password)
	// This is a placeholder implementation, replace with your actual logic
	if err := s.validator.Validate(req); err != nil {
		s.logger.Error(common.StandardMsgError(ctx, "reset password", err))
		return nil, status.Error(codes.InvalidArgument, "reset password fail by invalid data")
	}

	if req.NewPassword != req.ConfirmPassword {
		s.logger.Error(common.StandardMsgError(ctx, "reset password", errors.New("new password not match confirm password")))
		return nil, status.Error(codes.InvalidArgument, "new password not match confirm password")
	}

	// find token
	email, err := s.sessionStorage.GetSession(
		fmt.Sprintf("%v:%v", prefixForgotPasswordSessionStorage, req.ResetToken),
	)
	if err != nil {
		s.logger.Error(common.StandardMsgError(ctx, "reset password", err))
		return nil, status.Error(codes.Internal, "internal server")
	}

	// get user

	user, err := s.userRepo.FindOneByEmail(email)
	if err != nil {
		s.logger.Error(common.StandardMsgError(ctx, "reset password", err))
		return nil, status.Error(codes.Internal, "internal server")
	}

	hashPw, err := argon.EncodePassword(req.NewPassword)
	if err != nil {
		s.logger.Error(common.StandardMsgError(ctx, "reset password", err))
		return nil, status.Error(codes.Internal, "internal server")
	}

	_, err = s.userRepo.UpdateOneById(user.ID, map[string]interface{}{
		"password": hashPw,
	})
	if err != nil {
		s.logger.Error(common.StandardMsgError(ctx, "reset password", err))
		return nil, status.Error(codes.Internal, "internal server")
	}

	return &emptypb.Empty{}, nil
}

// SignUp implements the SignUp RPC method
func (s *AuthService) SignUp(ctx context.Context, req *pb.SignUpRequest) (*emptypb.Empty, error) {
	// Implement your signup logic here (e.g., create user account)
	// This is a placeholder implementation, replace with your actual logic
	s.logger.Infor(common.StandardMsgInfor(ctx, "sign up", "email:"+req.Email))

	if err := s.validator.Validate(req); err != nil {
		s.logger.Error(common.StandardMsgError(ctx, "sign up", err))
		return nil, status.Error(codes.InvalidArgument, "sign up fail by invalid data")
	}

	// parse to map
	mapUser, err := convert.StructProtoToMap(req)
	if err != nil {
		s.logger.Error(common.StandardMsgError(ctx, "sign up", err))
		return nil, status.Error(codes.Internal, "internal server")
	}

	if _, err := s.userRepo.CreateUser(mapUser); err != nil {
		s.logger.Error(common.StandardMsgError(ctx, "sign up", err))
		return nil, status.Error(codes.Internal, "internal server")
	}

	s.logger.Infor(
		common.StandardMsgInfor(ctx, "sign up", "success full email"+req.Email),
	)

	return &emptypb.Empty{}, nil
}
