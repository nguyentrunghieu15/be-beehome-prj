package provider

import (
	proapi "github.com/nguyentrunghieu15/be-beehome-prj/api/pro-api"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/validator"
)

func SetRules(v validator.IValidator) {
	// Add validation rules for your new structs here
	SetRulesOfFindProByIdRequest(v)
	SetRulesOfDeleteProByIdRequest(v)
	SetRulesOfSignUpProRequest(v)
	SetRulesOfUpdateProRequest(v)
	SetRulesOfAddPaymentMethodProRequest(v)
	SetRulesOfReplyReviewProRequest(v)
	SetRulesOfReviewProRequest(v)
	SetRulesOfAddServiceProRequest(v)
	SetRulesOfAddSocialMediaProRequest(v)
}

func SetRulesOfFindProByIdRequest(v validator.IValidator) {
	validationRules := map[string]string{
		"Id": "required,uuid",
	}
	v.RegisterRules(validationRules, &proapi.FindProByIdRequest{})
}

func SetRulesOfDeleteProByIdRequest(v validator.IValidator) {
	validationRules := map[string]string{
		"Id": "required,uuid",
	}
	v.RegisterRules(validationRules, &proapi.DeleteProByIdRequest{})
}

func SetRulesOfSignUpProRequest(v validator.IValidator) {
	validationRules := map[string]string{
		"Name":         "required",
		"Introduction": "required",
		"Years":        "gte=0", // Assuming years cannot be negative
		"PostalCode":   "required,postcode_iso3166_alpha2",
	}
	v.RegisterRules(validationRules, &proapi.SignUpProRequest{})
}

func SetRulesOfUpdateProRequest(v validator.IValidator) {
	validationRules := map[string]string{
		"Name":         "omitempty,required",
		"Introduction": "omitempty,required",
		"Years":        "omitempty,required,gte=0", // Assuming years cannot be negative
		"PostalCode":   "omitempty,required,postcode_iso3166_alpha2",
	}
	v.RegisterRules(validationRules, &proapi.UpdateProRequest{})
}

func SetRulesOfAddPaymentMethodProRequest(v validator.IValidator) {
	validationRules := map[string]string{
		"Name": "required",
	}
	v.RegisterRules(validationRules, &proapi.AddPaymentMethodProRequest{})
}

func SetRulesOfReplyReviewProRequest(v validator.IValidator) {
	validationRules := map[string]string{
		"ReviewId": "required,uuid",
		"Reply":    "required",
	}
	v.RegisterRules(validationRules, &proapi.ReplyReviewProRequest{})
}

func SetRulesOfReviewProRequest(v validator.IValidator) {
	validationRules := map[string]string{
		"ProviderId": "required,uuid",
		"Rating":     "required,gte=1,lte=5", // Assuming rating is between 1 and 5
		"Comment":    "omitempty",
	}
	v.RegisterRules(validationRules, &proapi.ReviewProRequest{})
}

func SetRulesOfAddServiceProRequest(v validator.IValidator) {
	validationRules := map[string]string{
		"ServiceId": "required,uuid",
	}
	v.RegisterRules(validationRules, &proapi.AddServiceProRequest{})
}

func SetRulesOfAddSocialMediaProRequest(v validator.IValidator) {
	validationRules := map[string]string{
		"Name":       "required",
		"Link":       "required,url", // Assuming link should be a valid url
		"ProviderId": "required,uuid",
	}
	v.RegisterRules(validationRules, &proapi.AddSocialMediaProRequest{})
}
