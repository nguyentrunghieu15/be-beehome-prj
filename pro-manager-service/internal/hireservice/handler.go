package hireservice

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	addressapi "github.com/nguyentrunghieu15/be-beehome-prj/api/address-api"
	proapi "github.com/nguyentrunghieu15/be-beehome-prj/api/pro-api"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/convert"
	communication "github.com/nguyentrunghieu15/be-beehome-prj/pro-manager-service/internal/comunitication"
	"github.com/nguyentrunghieu15/be-beehome-prj/pro-manager-service/mapper"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *HireService) FindAllHire(
	ctx context.Context,
	req *proapi.FindAllHireRequest,
) (*proapi.FindAllHireResponse, error) {
	// Validate the request
	if err := s.validator.Validate(req); err != nil {
		s.logger.Error(fmt.Sprintf("failed to validate FindAllHire request: %v", err))
		return nil, err
	}

	// Convert request to map for filtering
	filters, err := convert.StructProtoToMap(req)
	if err != nil {
		s.logger.Error(fmt.Sprintf("failed to convert request to map: %v", err))
		return nil, err
	}

	// Use hireRepo to fetch all hires based on filters
	hires, err := s.hireRepo.FindAll(filters)
	if err != nil {
		s.logger.Error(fmt.Sprintf("failed to fetch hires: %v", err))
		return nil, err
	}

	return &proapi.FindAllHireResponse{Hires: mapper.MapToListHireInfors(hires)}, nil
}

func (s *HireService) CreateHire(
	ctx context.Context,
	req *proapi.CreateHireRequest,
) (*proapi.CreateHireResponse, error) {
	// Validate request
	if err := s.validator.Validate(req); err != nil {
		return nil, err
	}

	mapHire, err := convert.StructProtoToMap(req)
	if err != nil {
		return nil, err
	}

	splitedAddress := strings.Split(req.Address, ",")
	if len(splitedAddress) != 3 {
		return nil, errors.New("Địa chỉ không đúng định dạng")
	}
	// check address
	isValidAddress, err := s.addressClient.CheckExistAddress(context.Background(), &addressapi.CheckExistAddressRequest{
		Address: &addressapi.Address{
			WardFullName:     strings.Trim(splitedAddress[0], " "),
			DistrictFullName: strings.Trim(splitedAddress[1], " "),
			ProvinceFullName: strings.Trim(splitedAddress[2], " "),
		},
	})

	if !isValidAddress {
		return nil, errors.New("Không tìm thấy địa chỉ")
	}

	userId := uuid.MustParse(ctx.Value("user_id").(string))
	mapHire["user_id"] = userId

	// Use hireRepo to create a new hire
	hire, err := s.hireRepo.CreateHire(mapHire)
	if err != nil {
		s.logger.Error(fmt.Sprintf("failed to create hire: %v", err))
		return nil, err
	}

	tranferMsg, err := json.Marshal(map[string]interface{}{
		"type":        "create",
		"provider_id": hire.ProviderId,
		"user_id":     userId.String(),
		"hire_id":     hire.ID.String(),
	})
	if err != nil {
		return nil, err
	}
	communication.HireResourceKafka.WriteMessages(
		context.Background(),
		kafka.Message{
			Value: tranferMsg,
		},
	)

	// Convert internal hire to CreateHireResponse format
	return &proapi.CreateHireResponse{Hire: mapper.MapToHire(hire)}, nil
}

func (s *HireService) UpdateStatusHire(
	ctx context.Context,
	req *proapi.UpdateStatusHireRequest,
) (*proapi.UpdateStatusHireResponse, error) {
	// Validate request
	if err := s.validator.Validate(req); err != nil {
		s.logger.Error(fmt.Sprintf("validation error for update status hire request: %v", err))
		return nil, err
	}

	// Parse Hire ID
	hireID := uuid.MustParse(req.GetHireId())

	// Update params
	updateParams := map[string]interface{}{"status": req.NewStatus}

	// Update Hire in database
	updatedHire, err := s.hireRepo.UpdateHireById(hireID, updateParams)
	if err != nil {
		s.logger.Error(fmt.Sprintf("failed to update hire status: %w", err))
		return nil, err
	}

	// Convert and return response (optional)
	// You can optionally convert the updatedHire object to a gRPC message format
	// using convertHireToResponse if needed for the response.

	return &proapi.UpdateStatusHireResponse{Hire: mapper.MapToHire(updatedHire)}, nil // Empty response for now
}

func (s *HireService) DeleteHire(ctx context.Context, req *proapi.DeleteHireRequest) (*emptypb.Empty, error) {
	// Validate request
	if err := s.validator.Validate(req); err != nil {
		s.logger.Error(fmt.Sprintf("validation error for update status hire request: %v", err))
		return nil, err
	}
	// Parse Hire ID
	hireID := uuid.MustParse(req.GetHireId())

	// Delete Hire from database
	err := s.hireRepo.DeleteHire(hireID)
	if err != nil {
		s.logger.Error(fmt.Sprintf("failed to delete hire: %w", err))
		return nil, err
	}

	tranferMsg, err := json.Marshal(map[string]interface{}{
		"type":    "delete",
		"hire_id": hireID.String(),
	})
	if err != nil {
		return nil, err
	}
	communication.HireResourceKafka.WriteMessages(
		context.Background(),
		kafka.Message{
			Value: tranferMsg,
		},
	)

	// Return empty response on successful deletion
	return &emptypb.Empty{}, nil
}
