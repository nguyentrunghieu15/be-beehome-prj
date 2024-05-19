package auth

import (
	authapi "github.com/nguyentrunghieu15/be-beehome-prj/api/auth-api"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/validator"
)

func SetRules(v validator.IValidator) {
	SetRulesOfLoginRequest(v)
}

func SetRulesOfLoginRequest(v validator.IValidator) {
	validationRules := map[string]string{
		"Email":    "required,email",
		"Password": "required,min=8",
	}
	v.RegisterRules(validationRules, &authapi.LoginRequest{})
}
