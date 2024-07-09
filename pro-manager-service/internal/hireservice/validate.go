package hireservice

import (
	"time"

	"github.com/go-playground/validator/v10"
	proapi "github.com/nguyentrunghieu15/be-beehome-prj/api/pro-api"
	validatorwrapper "github.com/nguyentrunghieu15/be-beehome-prj/internal/validator"
)

func SetRules(v validatorwrapper.IValidator) {
	// Register custom validators
	v.RegisterValidation("futuretime", FutureTime)
	v.RegisterValidation("timeafter", WorkTimeToAfterWorkTimeFrom)
	// Add validation rules for your new structs here
	SetRulesOfCreateHireRequest(v)       // New Rule
	SetRulesOfUpdateStatusHireRequest(v) // New Rule
	SetRulesOfDeleteHireRequest(v)       // New Rule
	SetRulesOfFindAllHireRequest(v)
}

func FutureTime(fl validator.FieldLevel) bool {
	if date, ok := fl.Field().Interface().(string); ok {
		layout := "2006-01-02T15:04:05Z" // adjust layout as per your time format
		t, err := time.Parse(layout, date)
		if err != nil {
			return false
		}
		return t.After(time.Now())
	}
	return false
}

func WorkTimeToAfterWorkTimeFrom(fl validator.FieldLevel) bool {
	workTimeTo, ok1 := fl.Field().Interface().(string)
	workTimeFrom, ok2 := fl.Parent().FieldByName("WorkTimeFrom").Interface().(string)
	if ok1 && ok2 {
		layout := "2006-01-02T15:04:05Z" // adjust layout as per your time format
		tTo, errTo := time.Parse(layout, workTimeTo)
		tFrom, errFrom := time.Parse(layout, workTimeFrom)
		if errTo != nil || errFrom != nil {
			return false
		}
		return tTo.After(tFrom)
	}
	return false
}

func SetRulesOfCreateHireRequest(v validatorwrapper.IValidator) {
	validationRules := map[string]string{
		"ProviderId":      "required,uuid",
		"ServiceId":       "required,uuid",
		"WorkTimeFrom":    "required,futuretime",
		"WorkTimeTo":      "required,timeafter",
		"Status":          "omitempty,oneof=pendding starting finished review cancel",
		"PaymentMethodId": "omitempty,uuid",
		"Issue":           "required",
		"Address":         "required",
		"FullAddress":     "required",
	}
	v.RegisterRules(validationRules, &proapi.CreateHireRequest{})
}

func SetRulesOfUpdateStatusHireRequest(v validatorwrapper.IValidator) {
	validationRules := map[string]string{
		"HireId":    "required,uuid",
		"NewStatus": "required,oneof=pendding starting finished review cancel",
	}
	v.RegisterRules(validationRules, &proapi.UpdateStatusHireRequest{})
}

func SetRulesOfDeleteHireRequest(v validatorwrapper.IValidator) {
	validationRules := map[string]string{
		"HireId": "required,uuid",
	}
	v.RegisterRules(validationRules, &proapi.DeleteHireRequest{})
}

func SetRulesOfFindAllHireRequest(v validatorwrapper.IValidator) {
	validationRules := map[string]string{
		"UserId":     "omitempty,uuid",
		"ProviderId": "omitempty,uuid",
		"Status":     "omitempty,oneof=pendding starting finished review cancel",
	}
	v.RegisterRules(validationRules, &proapi.FindAllHireRequest{})
}
