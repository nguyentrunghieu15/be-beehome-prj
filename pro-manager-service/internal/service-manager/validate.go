package servicemanager

import (
	"time"

	"github.com/go-playground/validator/v10"
	proapi "github.com/nguyentrunghieu15/be-beehome-prj/api/pro-api"
	validatorwrapper "github.com/nguyentrunghieu15/be-beehome-prj/internal/validator"
)

func ValidDate(fl validator.FieldLevel) bool {
	if date, ok := fl.Field().Interface().(string); ok {
		layout := "2006-01-02T15:04:05Z" // Adjust layout as per your time format
		_, err := time.Parse(layout, date)
		return err == nil
	}
	return false
}

func SetRules(v validatorwrapper.IValidator) {
	// Register custom validators
	v.RegisterValidation("validdate", ValidDate)

	// Add validation rules for your new structs here
	SetRulesOfFulltextSearchServicesRequest(v)
	SetRulesOfFilterServices(v)
	SetRulesOfListServicesRequest(v)
	SetRulesOfGetServiceRequest(v)
	SetRulesOfCreateServiceRequest(v)
	SetRulesOfUpdateServiceRequest(v)
	SetRulesOfDeleteServiceRequest(v)
}

func SetRulesOfFulltextSearchServicesRequest(v validatorwrapper.IValidator) {
	validationRules := map[string]string{
		"Name": "omitempty",
	}
	v.RegisterRules(validationRules, &proapi.FulltextSearchServicesRequest{})
}

func SetRulesOfFilterServices(v validatorwrapper.IValidator) {
	validationRules := map[string]string{
		"Name":          "omitempty",
		"CreatedAtFrom": "omitempty,validdate",
		"CreatedAtTo":   "omitempty,validdate",
		"UpdatedAtFrom": "omitempty,validdate",
		"UpdatedAtTo":   "omitempty,validdate",
	}
	v.RegisterRules(validationRules, &proapi.FilterServices{})
}

func SetRulesOfListServicesRequest(v validatorwrapper.IValidator) {
	validationRules := map[string]string{
		"Filter.Name":          "omitempty",
		"Filter.CreatedAtFrom": "omitempty,validdate",
		"Filter.CreatedAtTo":   "omitempty,validdate",
		"Filter.UpdatedAtFrom": "omitempty,validdate",
		"Filter.UpdatedAtTo":   "omitempty,validdate",
		"Pagination.Limit":     "omitempty,gte=0",
		"Pagination.Page":      "omitempty,gte=0",
		"Pagination.PageSize":  "omitempty,gte=0",
		"Pagination.Sort":      "omitempty",
		"Pagination.SortBy":    "omitempty",
	}
	v.RegisterRules(validationRules, &proapi.ListServicesRequest{})
}

func SetRulesOfGetServiceRequest(v validatorwrapper.IValidator) {
	validationRules := map[string]string{
		"Id": "required,uuid",
	}
	v.RegisterRules(validationRules, &proapi.GetServiceRequest{})
}

func SetRulesOfCreateServiceRequest(v validatorwrapper.IValidator) {
	validationRules := map[string]string{
		"Name":           "required",
		"Detail":         "required",
		"GroupServiceId": "required,uuid",
	}
	v.RegisterRules(validationRules, &proapi.CreateServiceRequest{})
}

func SetRulesOfUpdateServiceRequest(v validatorwrapper.IValidator) {
	validationRules := map[string]string{
		"Id":     "required,uuid",
		"Name":   "required",
		"Detail": "required",
	}
	v.RegisterRules(validationRules, &proapi.UpdateServiceRequest{})
}

func SetRulesOfDeleteServiceRequest(v validatorwrapper.IValidator) {
	validationRules := map[string]string{
		"Id": "required,uuid",
	}
	v.RegisterRules(validationRules, &proapi.DeleteServiceRequest{})
}
