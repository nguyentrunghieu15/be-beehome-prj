package groupservicemanager

import (
	"time"

	"github.com/go-playground/validator/v10"
	proapi "github.com/nguyentrunghieu15/be-beehome-prj/api/pro-api"
	validatorwrappr "github.com/nguyentrunghieu15/be-beehome-prj/internal/validator"
)

func ValidDate(fl validator.FieldLevel) bool {
	if date, ok := fl.Field().Interface().(string); ok {
		layout := "2006-01-02" // Adjust layout as per your time format
		_, err := time.Parse(layout, date)
		return err == nil
	}
	return false
}

func SetRules(v validatorwrappr.IValidator) {
	// Register custom validators
	v.RegisterValidation("validdate", ValidDate)

	// Add validation rules for your new structs here
	SetRulesOfFilterGroupServices(v)
	SetRulesOfPagination(v)
	SetRulesOfListGroupServicesRequest(v)
	SetRulesOfFulltextSearchGroupServicesRequest(v)
	SetRulesOfGetGroupServiceRequest(v)
	SetRulesOfCreateGroupServiceRequest(v)
	SetRulesOfUpdateGroupServiceRequest(v)
	SetRulesOfDeleteGroupServiceRequest(v)
}

func SetRulesOfFilterGroupServices(v validatorwrappr.IValidator) {
	validationRules := map[string]string{
		"Name":          "omitempty",
		"CreatedAtFrom": "omitempty,validdate",
		"CreatedAtTo":   "omitempty,validdate",
		"UpdatedAtFrom": "omitempty,validdate",
		"UpdatedAtTo":   "omitempty,validdate",
	}
	v.RegisterRules(validationRules, &proapi.FilterGroupServices{})
}

func SetRulesOfPagination(v validatorwrappr.IValidator) {
	validationRules := map[string]string{
		"Limit":    "omitempty,gte=0",
		"Page":     "omitempty,gte=0",
		"PageSize": "omitempty,gte=0",
		"Sort":     "omitempty",
		"SortBy":   "omitempty",
	}
	v.RegisterRules(validationRules, &proapi.Pagination{})
}

func SetRulesOfListGroupServicesRequest(v validatorwrappr.IValidator) {
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
	v.RegisterRules(validationRules, &proapi.ListGroupServicesRequest{})
}

func SetRulesOfFulltextSearchGroupServicesRequest(v validatorwrappr.IValidator) {
	validationRules := map[string]string{
		"Name": "omitempty",
	}
	v.RegisterRules(validationRules, &proapi.FulltextSearchGroupServicesRequest{})
}

func SetRulesOfGetGroupServiceRequest(v validatorwrappr.IValidator) {
	validationRules := map[string]string{
		"Id": "required,uuid",
	}
	v.RegisterRules(validationRules, &proapi.GetGroupServiceRequest{})
}

func SetRulesOfCreateGroupServiceRequest(v validatorwrappr.IValidator) {
	validationRules := map[string]string{
		"Name":   "required",
		"Detail": "required",
	}
	v.RegisterRules(validationRules, &proapi.CreateGroupServiceRequest{})
}

func SetRulesOfUpdateGroupServiceRequest(v validatorwrappr.IValidator) {
	validationRules := map[string]string{
		"Id":     "required,uuid",
		"Name":   "required",
		"Detail": "required",
	}
	v.RegisterRules(validationRules, &proapi.UpdateGroupServiceRequest{})
}

func SetRulesOfDeleteGroupServiceRequest(v validatorwrappr.IValidator) {
	validationRules := map[string]string{
		"Id": "required,uuid",
	}
	v.RegisterRules(validationRules, &proapi.DeleteGroupServiceRequest{})
}
