package grpcaddress

import (
	"context"
	"errors"

	"github.com/nguyentrunghieu15/be-beehome-prj/address-service/internal/address"
	"github.com/nguyentrunghieu15/be-beehome-prj/address-service/internal/mapper"
	addressapi "github.com/nguyentrunghieu15/be-beehome-prj/api/address-api"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/logwrapper"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/validator"
)

type AddressService struct {
	addressapi.UnimplementedAddressServiceServer
	addressRepo address.IAddressRepo
	validator   validator.IValidator
	logger      logwrapper.ILoggerWrapper
}

func NewAddressService(addressRepo address.IAddressRepo,
	validator validator.IValidator,
	logger logwrapper.ILoggerWrapper,
) *AddressService {
	return &AddressService{addressRepo: addressRepo, validator: validator, logger: logger}

}

const (
	PROVINCE = "province"
	DISTRICT = "district"
	WARD     = "ward"
)

func (s *AddressService) CheckExistAddress(ctx context.Context,
	req *addressapi.CheckExistAddressRequest,
) (*addressapi.CheckExistAddressResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	if exsited, err := s.addressRepo.CheckAddressExist(*mapper.ConvertAddress(req.GetAddress())); err != nil {
		s.logger.Error(err.Error())
		return nil, err
	} else if exsited {
		return &addressapi.CheckExistAddressResponse{
			Address: req.Address,
		}, nil
	}
	return nil, nil
}

func (s *AddressService) FindAllAddress(ctx context.Context,
	req *addressapi.FindAllAddressRequest,
) (*addressapi.FindAllAddressResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	result, err := s.addressRepo.FindAllAddressByQuery(req.Query)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	return &addressapi.FindAllAddressResponse{
		Address: mapper.ConvertListDatasourceAddressToString(result),
	}, nil
}

func (s *AddressService) FindAllNameOfAddressUnit(ctx context.Context,
	req *addressapi.FindAllNameOfAddressUnitRequest,
) (*addressapi.FindAllNameOfAddressUnitResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}
	switch req.Type {
	case PROVINCE:
		results, err := s.addressRepo.GetAllProvince()
		if err != nil {
			s.logger.Error(err.Error())
			return nil, err
		}

		var resultInterfaces []interface{} = make([]interface{}, len(results))

		// Convert each Province struct to an interface
		for i, result := range results {
			resultInterfaces[i] = result
		}

		res, err := mapper.ToListAddressUnit(resultInterfaces)
		if err != nil {
			s.logger.Error(err.Error())
			return nil, err
		}
		return &addressapi.FindAllNameOfAddressUnitResponse{
			Unit: res,
		}, nil
	case DISTRICT:
		if req.Unit == nil {
			err := errors.New("Thiếu dữ liệu")
			s.logger.Error(err.Error())
			return nil, err
		}
		results, err := s.addressRepo.GetAllDistrictsByProvinceCode(req.Unit.Code)
		if err != nil {
			s.logger.Error(err.Error())
			return nil, err
		}
		var resultInterfaces []interface{} = make([]interface{}, len(results))

		// Convert each Province struct to an interface
		for i, result := range results {
			resultInterfaces[i] = result
		}

		res, err := mapper.ToListAddressUnit(resultInterfaces)
		if err != nil {
			s.logger.Error(err.Error())
			return nil, err
		}
		return &addressapi.FindAllNameOfAddressUnitResponse{
			Unit: res,
		}, nil
	case WARD:
		if req.Unit == nil {
			err := errors.New("Thiếu dữ liệu")
			s.logger.Error(err.Error())
			return nil, err
		}
		results, err := s.addressRepo.GetAllWardsByDistrictCode(req.Unit.Code)
		if err != nil {
			s.logger.Error(err.Error())
			return nil, err
		}
		var resultInterfaces []interface{} = make([]interface{}, len(results))

		// Convert each Province struct to an interface
		for i, result := range results {
			resultInterfaces[i] = result
		}

		res, err := mapper.ToListAddressUnit(resultInterfaces)
		if err != nil {
			s.logger.Error(err.Error())
			return nil, err
		}
		return &addressapi.FindAllNameOfAddressUnitResponse{
			Unit: res,
		}, nil
	}
	return nil, nil
}
