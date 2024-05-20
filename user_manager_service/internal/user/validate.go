package user

import (
	userapi "github.com/nguyentrunghieu15/be-beehome-prj/api/user-api"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/validator"
)

func SetRules(v validator.IValidator) {
	SetRulesOfBlockRequest(v)
	SetRulesOfCreateUserRequest(v)
	SetRulesOfDeleteUserRequest(v)
	SetRulesOfGetUserRequest(v)
	SetRulesOfListUsersRequest(v)
	SetRulesOfUpdateUserRequest(v)
}

func SetRulesOfBlockRequest(v validator.IValidator) {
	validationRules := map[string]string{
		"Id": "required,uuid", // Assuming Id should be a valid UUID
	}
	v.RegisterRules(validationRules, &userapi.BlockRequest{})
}

func SetRulesOfCreateUserRequest(v validator.IValidator) {
	validationRules := map[string]string{
		"Email":     "required,email",
		"Password":  "required,min=8",
		"Phone":     "required,min=9", // Adjust min length based on your phone number format
		"FirstName": "required",
		"LastName":  "required",
	}
	v.RegisterRules(validationRules, &userapi.CreateUserRequest{})
}

func SetRulesOfListUsersRequest(v validator.IValidator) {
	// You might not need validation for ListUsersRequest
	SetRulesOfPagination(v)
	SetRulesOfSort(v)
	// But you can add rules for pagination or sort fields if required
}

func SetRulesOfGetUserRequest(v validator.IValidator) {
	validationRules := map[string]string{
		"Id": "required,uuid", // Assuming Id should be a valid UUID
	}
	v.RegisterRules(validationRules, &userapi.GetUserRequest{})
}

func SetRulesOfUpdateUserRequest(v validator.IValidator) {
	validationRules := map[string]string{
		"Id":        "required,uuid",   // Assuming Id should be a valid UUID
		"Email":     "omitempty,email", // Email can be empty or a valid email
		"Phone":     "omitempty,min=9", // Phone can be empty or follow min length rule
		"FirstName": "omitempty",       // Other fields can be empty
		"LastName":  "omitempty",
		"Status":    "omitempty", // Assuming Status has predefined valid values, you can add a custom validation rule here
	}
	v.RegisterRules(validationRules, &userapi.UpdateUserRequest{})
}

func SetRulesOfDeleteUserRequest(v validator.IValidator) {
	validationRules := map[string]string{
		"Id": "required,uuid", // Assuming Id should be a valid UUID
	}
	v.RegisterRules(validationRules, &userapi.DeleteUserRequest{})
}

func SetRulesOfPagination(v validator.IValidator) {
	validationRules := map[string]string{
		"Pagination.PageSize":  "omitempty,gte=1", // optional, but must be greater than or equal to 1
		"Pagination.PageToken": "omitempty",       // optional
	}
	v.RegisterRules(validationRules, &userapi.Pagination{})
}

func SetRulesOfSort(v validator.IValidator) {
	validationRules := map[string]string{
		"Type":  "omitempty,required,oneof=asc desc", // must be either "asc" or "desc"
		"Field": "omitempty,required",                // field name cannot be empty
	}
	v.RegisterRules(validationRules, &userapi.Sort{})
}
