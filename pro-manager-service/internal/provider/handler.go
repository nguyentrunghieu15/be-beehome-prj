package provider

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
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/random"
	"github.com/nguyentrunghieu15/be-beehome-prj/pkg/jwt"
	communication "github.com/nguyentrunghieu15/be-beehome-prj/pro-manager-service/internal/comunitication"
	"github.com/nguyentrunghieu15/be-beehome-prj/pro-manager-service/mapper"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Find providers based on search criteria
func (s *ProviderService) FindPros(ctx context.Context, req *proapi.FindProsRequest) (*proapi.FindProsResponse, error) {
	// Implement logic to search providers using proRepo based on request criteria
	// ...
	// Validate the request
	if err := s.validator.Validate(req); err != nil {
		return nil, err
	}

	providers, err := s.proRepo.FindProviders(req)
	if err != nil {
		return nil, err
	}

	return &proapi.FindProsResponse{
		Providers: mapper.MapToProviderViewInfos(providers),
	}, nil // Replace with actual response population
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

	tranferMsg, err := json.Marshal(map[string]interface{}{
		"type":        "delete",
		"provider_id": req.GetId(),
		"user_id":     ctx.Value("user_id").(string),
	})
	if err != nil {
		return nil, err
	}
	communication.ProviderResourceKafka.WriteMessages(
		context.Background(),
		kafka.Message{
			Value: tranferMsg,
		},
	)

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

	tranferMsg, err := json.Marshal(map[string]interface{}{
		"type":        "create",
		"provider_id": createdPro.ID.String(),
		"user_id":     userId.String(),
	})
	if err != nil {
		return nil, err
	}
	communication.ProviderResourceKafka.WriteMessages(
		context.Background(),
		kafka.Message{
			Value: tranferMsg,
		},
	)
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

	if req.Address != nil {
		splitedAddress := strings.Split(*req.Address, ",")
		if len(splitedAddress) != 3 {
			return nil, errors.New("Địa chỉ không đúng định dạng")
		}
		// check address
		isValidAddress, _ := s.addressClient.CheckExistAddress(
			context.Background(),
			&addressapi.CheckExistAddressRequest{
				Address: &addressapi.Address{
					WardFullName:     strings.Trim(splitedAddress[0], " "),
					DistrictFullName: strings.Trim(splitedAddress[1], " "),
					ProvinceFullName: strings.Trim(splitedAddress[2], " "),
				},
			},
		)
		if !isValidAddress {
			return nil, errors.New("Không tìm thấy địa chỉ")
		}
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
	paymentMethodData["id"] = random.GenerateRandomUUID()

	// Create payment method using paymentRepo

	_, err := s.paymentRepo.CreatePaymentMethod(paymentMethodData)
	if err != nil {
		return nil, err
	}

	tranferMsg, err := json.Marshal(map[string]interface{}{
		"type":              "create",
		"provider_id":       providerID.String(),
		"payment_method_id": paymentMethodData["id"],
	})
	if err != nil {
		return nil, err
	}
	communication.PaymentMethodResourceKafka.WriteMessages(
		context.Background(),
		kafka.Message{
			Value: tranferMsg,
		},
	)

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
	review, err := s.reviewRepo.CreateReview(data)
	if err != nil {
		return nil, err
	}

	tranferMsg, err := json.Marshal(map[string]interface{}{
		"type":        "create",
		"review_id":   review.ID.String(),
		"provider_id": providerID.String(),
		"user_id":     userID.String(),
		"hire_id":     req.HireId,
	})

	if err != nil {
		return nil, err
	}
	communication.SocialMediaResourceKafka.WriteMessages(
		context.Background(),
		kafka.Message{
			Value: tranferMsg,
		},
	)

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
	servicesID := make([]uuid.UUID, 0)
	for _, id := range req.GetServicesId() {
		servicesID = append(servicesID, uuid.MustParse(id))
	}
	// Add service to provider using providerRepo
	err := s.proRepo.AddServicesForPro(providerID, servicesID...)
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
	providerID := uuid.MustParse(
		ctx.Value("provider_id").(string),
	) // Implement this function based on your context usage

	// Convert request data to map
	data, err := convert.StructProtoToMap(req)
	if err != nil {
		return nil, err
	}

	// Add provider ID to data
	data["provider_id"] = providerID

	// Create social media record using socialMediaRepo
	socialMedia, err := s.socialMediaRepo.CreateSocialMedia(data)
	if err != nil {
		return nil, err
	}

	tranferMsg, err := json.Marshal(map[string]interface{}{
		"type":            "create",
		"social_media_id": socialMedia.ID.String(),
		"provider_id":     providerID.String(),
	})

	if err != nil {
		return nil, err
	}
	communication.SocialMediaResourceKafka.WriteMessages(
		context.Background(),
		kafka.Message{
			Value: tranferMsg,
		},
	)

	// Return empty response (modify if needed)
	return &emptypb.Empty{}, nil
}

func (s *ProviderService) JoinAsProvider(
	ctx context.Context,
	req *proapi.JoinAsProviderRequest,
) (*proapi.JoinAsProviderResponse, error) {
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

func (s *ProviderService) GetProviderProfile(
	ctx context.Context,
	req *emptypb.Empty,
) (*proapi.ProviderProfileResponse, error) {
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

func (s *ProviderService) GetAllServiceOfProvider(
	ctx context.Context,
	req *proapi.GetAllServiceOfProviderRequest,
) (*proapi.GetAllServiceOfProviderResponse, error) {
	// Validate the request
	if err := s.validator.Validate(req); err != nil {
		return nil, err
	}

	id := uuid.MustParse(req.GetId())
	services, err := s.proRepo.GetAllServicesOfProvider(id)
	if err != nil {
		return nil, err
	}
	return &proapi.GetAllServiceOfProviderResponse{Services: mapper.MapToServices(services)}, nil
}

// Add a service offered by a provider
func (s *ProviderService) DeleteServicePro(
	ctx context.Context,
	req *proapi.DeleteServiceProRequest,
) (*emptypb.Empty, error) {
	// Validate the request
	if err := s.validator.Validate(req); err != nil {
		return nil, err
	}

	// Extract provider ID from context (assuming context carries provider ID)
	providerID := uuid.MustParse(ctx.Value("provider_id").(string))

	// Convert service ID from string to uuid.UUID
	servicesID := make([]uuid.UUID, 0)
	for _, id := range req.GetServicesId() {
		servicesID = append(servicesID, uuid.MustParse(id))
	}
	// Add service to provider using providerRepo
	err := s.proRepo.RemoveServicesOfPro(providerID, servicesID...)
	if err != nil {
		return nil, err
	}

	// Return empty response (modify if needed)
	return &emptypb.Empty{}, nil
}

func (s *ProviderService) GetAllReviewsOfProvider(
	ctx context.Context,
	req *proapi.GetAllReviewOfProviderRequest,
) (*proapi.GetAllReviewOfProviderResponse, error) {
	// Validate the request
	if err := s.validator.Validate(req); err != nil {
		return nil, err
	}

	// Extract provider ID from context (assuming context carries provider ID)
	providerID := uuid.MustParse(req.GetId())

	reviews, err := s.reviewRepo.FindReviewsByProviderId(providerID)
	if err != nil {
		return nil, err
	}

	return &proapi.GetAllReviewOfProviderResponse{
		Reviews: mapper.MapToReviews(reviews),
	}, nil
}

func (s *ProviderService) UpdateSocialMediaPro(
	ctx context.Context,
	req *proapi.UpdateSocialMediaProRequest,
) (*emptypb.Empty, error) {
	// Validate the request
	if err := s.validator.Validate(req); err != nil {
		return nil, err
	}

	mapData, err := convert.StructProtoToMap(req)
	if err != nil {
		return nil, err
	}

	delete(mapData, "id")

	_, err = s.socialMediaRepo.UpdateOneById(uuid.MustParse(req.Id), mapData)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// Add social media information for a professional
func (s *ProviderService) DeleteSocialMediaPro(
	ctx context.Context,
	req *proapi.DeleteSocialMediaProRequest,
) (*emptypb.Empty, error) {
	// Validate the request
	if err := s.validator.Validate(req); err != nil {
		return nil, err
	}
	id := uuid.MustParse(req.Id)
	err := s.socialMediaRepo.DeleteOneById(id)
	if err != nil {
		return nil, err
	}

	tranferMsg, err := json.Marshal(map[string]interface{}{
		"type":            "delete",
		"social_media_id": id.String(),
	})

	if err != nil {
		return nil, err
	}
	communication.SocialMediaResourceKafka.WriteMessages(
		context.Background(),
		kafka.Message{
			Value: tranferMsg,
		},
	)

	return &emptypb.Empty{}, nil
}
