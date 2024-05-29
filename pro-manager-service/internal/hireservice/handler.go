package hireservice

import (
	"context"
	"fmt"

	proapi "github.com/nguyentrunghieu15/be-beehome-prj/api/pro-api"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/convert"
	"github.com/nguyentrunghieu15/be-beehome-prj/pro-manager-service/mapper.go"
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

	return &proapi.FindAllHireResponse{Hires: mapper.MapToListHire(hires)}, nil
}

func (s *HireService) CreateHire(
	ctx context.Context,
	req *proapi.CreateHireRequest,
) (*proapi.CreateHireResponse, error) {
	// Validate request using validator
	err := s.validator.Validate(req)
	if err != nil {
		s.logger.Error(fmt.Sprintf("validation error for create hire request: %v", err))
		return nil, err
	}

	// Use hireRepo to create a new hire
	hire, err := s.hireRepo.Create(ctx, req)
	if err != nil {
		s.logger.Error(fmt.Sprintf("failed to create hire: %v", err))
		return nil, err
	}

	// Convert internal hire to CreateHireResponse format
	return &CreateHireResponse{Hire: convertHireToResponse(hire)}, nil
}

func (s *HireService) UpdateStatusHire(
	ctx context.Context,
	req *proapi.UpdateStatusHireRequest,
) (*proapi.UpdateStatusHireResponse, error) {
	// Validate request using validator
	err := s.validator.Validate(req)
	if err != nil {
		s.logger.Error(fmt.Sprintf("validation error for update status hire request: %v", err))
		return nil, err
	}

	// Use hireRepo to update hire status
	err = s.hireRepo.UpdateStatus(ctx, req.HireId, req.Status)
	if err != nil {
		s.logger.Error(fmt.Sprintf("failed to update hire status: %v", err))
		return nil, err
	}

	// No response object needed for update, return empty response
	return &UpdateStatusHireResponse{}, nil
}

func (s *HireService) DeleteHire(ctx context.Context, req *proapi.DeleteHireRequest) (*emptypb.Empty, error) {
	// Use hireRepo to delete hire
	err := s.hireRepo.Delete(ctx, req.HireId)
	if err != nil {
		s.logger.Error(fmt.Sprintf("failed to delete hire: %v", err))
		return nil, err
	}

	// Return empty response on successful deletion
	return &emptypb.Empty{}, nil
}
