package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	proapi "github.com/nguyentrunghieu15/be-beehome-prj/api/pro-api"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/convert"
	"github.com/nguyentrunghieu15/be-beehome-prj/pkg/jwt"
	"github.com/nguyentrunghieu15/be-beehome-prj/pro-manager-service/mapper.go"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Find providers based on search criteria
func (s *ProviderService) FindPros(ctx context.Context, req *proapi.FindProsRequest) (*proapi.FindProsResponse, error) {
	// Implement logic to search providers using proRepo based on request criteria
	// ...
	return &proapi.FindProsResponse{}, nil // Replace with actual response population
}

// Find pro by ID
func (s *ProviderService) FindProById(
	ctx context.Context,
	req *proapi.FindProByIdRequest,
) (*proapi.FindProByIdResponse, error) {
	// Validate the request
	if err := s.validator.Validate(req); err != nil {
		return nil, err
	}

	// Use proRepo to find the provider by ID
	provider, err := s.proRepo.FindOneById(uuid.MustParse(req.GetId()))
	if err != nil {
		return nil, err
	}

	// Map Provider struct to ProviderInfo struct
	providerInfo := mapper.MapProviderToInfo(provider)

	// Build and return response
	return &proapi.FindProByIdResponse{
		Provider: providerInfo,
	}, nil
}

// Delete pro by ID
func (s *ProviderService) DeleteProById(ctx context.Context, req *proapi.DeleteProByIdRequest) (*emptypb.Empty, error) {
	// Validate the request
	if err := s.validator.Validate(req); err != nil {
		return nil, err
	}

	// Use proRepo to delete the provider by ID
	err := s.proRepo.DeleteOneById(uuid.MustParse(req.GetId()))
	if err != nil {
		return nil, err
	}

	// Delete successful response
	return nil, nil
}

// Sign up as a professional
func (s *ProviderService) SignUpPro(ctx context.Context, req *proapi.SignUpProRequest) (*proapi.ProviderInfo, error) {
	// Validate the request
	if err := s.validator.Validate(req); err != nil {
		return nil, err
	}

	// Convert request to map
	data, err := convert.StructProtoToMap(req)
	if err != nil {
		return nil, err
	}
	delete(data, "postal_code")

	postalCode, err := s.postalCodeRepo.FindPostalCodesByZipcode(req.PostalCode)
	if err != nil {
		return nil, err
	}
	if len(postalCode) == 0 {
		return nil, errors.New("not found zip code")
	}
	data["postal_code_id"] = postalCode[0].ID

	userId := uuid.MustParse(ctx.Value("user_id").(string))
	data["user_id"] = userId

	providerAlready, _ := s.proRepo.FindOneByUserId(userId)
	if providerAlready != nil {
		return nil, errors.New("Provider exist")
	}

	// Create a new provider record
	createdPro, err := s.proRepo.CreateProvider(data)
	if err != nil {
		return nil, err
	}

	// Potentially interact with paymentRepo for payment methods (not implemented here)
	// You might need to add logic to handle payment methods based on your requirements

	// Return empty response (modify if needed)
	return mapper.MapProviderToInfo(createdPro), nil
}

// Update information of a professional
func (s *ProviderService) UpdatePro(ctx context.Context, req *proapi.UpdateProRequest) (*proapi.ProviderInfo, error) {
	// Validate the request
	if err := s.validator.Validate(req); err != nil {
		return nil, err
	}

	// Convert request to map (optional, if GORM tags are not used)
	updateData, err := convert.StructProtoToMap(req)
	if err != nil {
		return nil, err
	}

	if _, ok := updateData["postal_code"]; ok {
		postalCode, err := s.postalCodeRepo.FindPostalCodesByZipcode(req.GetPostalCode())
		if err != nil {
			return nil, err
		}
		if len(postalCode) == 0 {
			return nil, errors.New("not found zip code")
		}
		delete(updateData, "postal_code")
		updateData["postal_code_id"] = postalCode[0].ID
	}

	// Update provider using GORM with associations (recommended)
	updatedPro, err := s.proRepo.UpdateOneById(uuid.MustParse(req.Id), updateData)
	if err != nil {
		return nil, err
	}

	// Return empty response (modify if needed)
	return mapper.MapProviderToInfo(updatedPro), nil
}

// Add a payment method for a provider
func (s *ProviderService) AddPaymentMethodPro(
	ctx context.Context,
	req *proapi.AddPaymentMethodProRequest,
) (*emptypb.Empty, error) {
	// Validate the request
	if err := s.validator.Validate(req); err != nil {
		return nil, err
	}

	// Extract provider ID from context (assuming context carries provider ID)
	providerID := uuid.MustParse(
		ctx.Value("provider_id").(string),
	) // Implement this function based on your context usage

	// Create payment method data
	paymentMethodData := map[string]interface{}{
		"name":        req.GetName(),
		"provider_id": providerID, // Use extracted provider ID
	}

	// Create payment method using paymentRepo
	_, err := s.paymentRepo.CreatePaymentMethod(paymentMethodData)
	if err != nil {
		return nil, err
	}

	// Return empty response (modify if needed)
	return &emptypb.Empty{}, nil
}

// Reply to a review as a professional
func (s *ProviderService) ReplyReviewPro(
	ctx context.Context,
	req *proapi.ReplyReviewProRequest,
) (*emptypb.Empty, error) {
	// Validate the request
	if err := s.validator.Validate(req); err != nil {
		return nil, err
	}

	// Extract review ID from request
	reviewID, err := uuid.Parse(req.GetReviewId())
	if err != nil {
		return nil, fmt.Errorf("invalid review ID format: %w", err)
	}

	// Convert request data to map
	data, err := convert.StructProtoToMap(req)
	if err != nil {
		return nil, err
	}
	delete(data, "review_id")

	// Update review record with reply using reviewRepo
	_, err = s.reviewRepo.UpdateOneById(reviewID, data)
	if err != nil {
		return nil, err
	}

	// Return empty response (modify if needed)
	return &emptypb.Empty{}, nil
}

// Review a professional as a user
func (s *ProviderService) ReviewPro(ctx context.Context, req *proapi.ReviewProRequest) (*emptypb.Empty, error) {
	// Validate the request
	if err := s.validator.Validate(req); err != nil {
		return nil, err
	}

	// Extract user ID from context (assuming context carries user ID)
	userID := uuid.MustParse(ctx.Value("user_id").(string))

	// Extract provider ID from request
	providerID, err := uuid.Parse(req.GetProviderId())
	if err != nil {
		return nil, fmt.Errorf("invalid provider ID format: %w", err)
	}

	// Convert request data to map
	data, err := convert.StructProtoToMap(req)
	if err != nil {
		return nil, err
	}

	// Add user ID and provider ID to data
	data["user_id"] = userID
	data["provider_id"] = providerID

	// Create review record using reviewRepo
	_, err = s.reviewRepo.CreateReview(data)
	if err != nil {
		return nil, err
	}

	// Return empty response (modify if needed)
	return &emptypb.Empty{}, nil
}

// Add a service offered by a provider
func (s *ProviderService) AddServicePro(ctx context.Context, req *proapi.AddServiceProRequest) (*emptypb.Empty, error) {
	// Validate the request
	if err := s.validator.Validate(req); err != nil {
		return nil, err
	}

	// Extract provider ID from context (assuming context carries provider ID)
	providerID := uuid.MustParse(ctx.Value("provider_id").(string))

	// Convert service ID from string to uuid.UUID
	serviceID := uuid.MustParse(req.GetServiceId())

	// Add service to provider using providerRepo
	err := s.proRepo.AddServicesForPro(providerID, serviceID)
	if err != nil {
		return nil, err
	}

	// Return empty response (modify if needed)
	return &emptypb.Empty{}, nil
}

// Add social media information for a provider
func (s *ProviderService) AddSocialMediaPro(
	ctx context.Context,
	req *proapi.AddSocialMediaProRequest,
) (*emptypb.Empty, error) {
	// Validate the request
	if err := s.validator.Validate(req); err != nil {
		return nil, err
	}

	// Extract provider ID from context (assuming context carries provider ID)
	providerID := uuid.MustParse(ctx.Value("provider_id").(string))

	// Convert request data to map
	data, err := convert.StructProtoToMap(req)
	if err != nil {
		return nil, err
	}

	// Add provider ID to data
	data["provider_id"] = providerID

	// Create social media record using socialMediaRepo
	_, err = s.socialMediaRepo.CreateSocialMedia(data)
	if err != nil {
		return nil, err
	}

	// Return empty response (modify if needed)
	return &emptypb.Empty{}, nil
}

func (s *ProviderService) JoinAsProvider(ctx context.Context, req *proapi.JoinAsProviderRequest) (*proapi.JoinAsProviderResponse, error) {
	// Validate the request
	if err := s.validator.Validate(req); err != nil {
		return nil, err
	}

	// Extract user ID from context (assuming context carries user ID)
	userID := uuid.MustParse(ctx.Value("user_id").(string))

	pro, err := s.proRepo.FindOneByUserId(userID)
	if err != nil {
		return nil, err
	}

	proToken, err := s.jwtTokenizer.GenerateToken(pro.ID.String(), jwt.DefaultAccessTokenConfigure)
	if err != nil {
		return nil, err
	}
	return &proapi.JoinAsProviderResponse{ProviderToken: proToken}, nil
}

func (s *ProviderService) GetProviderProfile(ctx context.Context, req *emptypb.Empty) (*proapi.ProviderProfileResponse, error) {
	// Validate the request
	if err := s.validator.Validate(req); err != nil {
		return nil, err
	}

	// Extract provider ID from context (assuming context carries provider ID)
	providerID := uuid.MustParse(ctx.Value("provider_id").(string))

	provider, err := s.proRepo.FindOneById(providerID)
	if err != nil {
		return nil, err
	}
	return &proapi.ProviderProfileResponse{Provider: mapper.MapProviderToInfo(provider)}, nil
}
