package profiles

import (
	userapi "github.com/nguyentrunghieu15/be-beehome-prj/api/user-api"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/validator"
)

func SetRules(v validator.IValidator) {
	SetRulesOfAddCardRequest(v)
	SetRulesOfCard(v)
	SetRulesOfChangeEmailRequest(v)
}

func SetRulesOfAddCardRequest(v validator.IValidator) {
	validationRules := map[string]string{
		"Card": "required", // Assuming Card field is required and has its own validation rules
	}
	v.RegisterRules(validationRules, &userapi.AddCardRequest{})
}

func SetRulesOfChangeEmailRequest(v validator.IValidator) {
	validationRules := map[string]string{
		"Email": "required,email",
	}
	v.RegisterRules(validationRules, &userapi.ChangeEmailRequest{})
}
func SetRulesOfCard(v validator.IValidator) {
	validationRules := map[string]string{
		"CardNumber": "required,len=16", // Adjust length based on your card number format
		"OwnerName":  "required",
		"BankName":   "required",
	}
	v.RegisterRules(validationRules, &userapi.Card{})
}
