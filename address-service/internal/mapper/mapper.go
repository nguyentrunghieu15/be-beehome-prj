package mapper

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/nguyentrunghieu15/be-beehome-prj/address-service/internal/datasource"
	addressapi "github.com/nguyentrunghieu15/be-beehome-prj/api/address-api"
)

func ConvertAddress(apiAddress *addressapi.Address) *datasource.Address {
	// Create a new datasource.Address object
	dsAddress := &datasource.Address{
		WardFullName:     apiAddress.GetWardFullName(),
		DistrictFullName: apiAddress.GetDistrictFullName(),
		ProvinceFullName: apiAddress.GetProvinceFullName(),
	}
	return dsAddress
}

// ConvertDatasourceAddress converts a datasource.Address to an addressapi.Address
func ConvertDatasourceAddressToString(dsAddress datasource.Address) string {
	// Create a new addressapi.Address object

	return fmt.Sprintf(
		"%s, %s, %s",
		dsAddress.WardFullName,
		dsAddress.DistrictFullName,
		dsAddress.ProvinceFullName,
	)
}

func ConvertListDatasourceAddressToString(dsAddress []datasource.Address) []string {
	// Create a new addressapi.Address object
	result := make([]string, 0)
	for _, v := range dsAddress {
		result = append(result, ConvertDatasourceAddressToString(v))
	}
	return result
}

func GetCodeAndFullName(data interface{}) (string, string, error) {
	v := reflect.ValueOf(data)
	t := v.Type()
	if t.Kind() != reflect.Struct {
		return "", "", errors.New("data must be a struct")
	}

	codeValue := ""
	fullNameValue := ""

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldName := field.Name

		if fieldName == "Code" {
			fieldValue := v.Field(i)
			if !fieldValue.CanInterface() {
				return "", "", errors.New("Code field is not accessible")
			}
			codeValue = fieldValue.Interface().(string)
		} else if fieldName == "FullName" {
			fieldValue := v.Field(i)
			if !fieldValue.CanInterface() {
				return "", "", errors.New("FullName field is not accessible")
			}
			fullNameValue = fieldValue.Interface().(string)
		}
	}

	return codeValue, fullNameValue, nil
}

func ToAddressUnit(data interface{}) (*addressapi.AddressUnit, error) {
	code, fullName, err := GetCodeAndFullName(data)
	if err != nil {
		return nil, err
	}
	return &addressapi.AddressUnit{
		Code: code,
		Name: fullName,
	}, nil
}

func ToListAddressUnit(data []interface{}) ([]*addressapi.AddressUnit, error) {
	result := []*addressapi.AddressUnit{}

	for _, v := range data {
		temp, err := ToAddressUnit(v)
		if err != nil {
			return nil, err
		}
		result = append(result, temp)
	}
	return result, nil
}
