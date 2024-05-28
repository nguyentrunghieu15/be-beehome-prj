package auth

import (
	authapi "github.com/nguyentrunghieu15/be-beehome-prj/api/auth-api"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/validator"
)

func SetRules(v validator.IValidator) {
	SetRulesOfLoginRequest(v)
	SetRulesOfForgetPasswordRequest(v)
	SetRulesOfRefreshTokenRequest(v)
	SetRulesOfResetPasswordRequest(v)
	SetRulesOfSignUpRequest(v)
}

func SetRulesOfLoginRequest(v validator.IValidator) {
	validationRules := map[string]string{
		"Email":    "required,email",
		"Password": "required,min=8",
	}
	v.RegisterRules(validationRules, &authapi.LoginRequest{})
}

func SetRulesOfRefreshTokenRequest(v validator.IValidator) {
	validationRules := map[string]string{
		"RefreshToken": "required",
	}
	v.RegisterRules(validationRules, &authapi.RefreshTokenRequest{})
}

func SetRulesOfForgetPasswordRequest(v validator.IValidator) {
	validationRules := map[string]string{
		"Email": "required,email",
	}
	v.RegisterRules(validationRules, &authapi.ForgotPasswordRequest{})
}

func SetRulesOfResetPasswordRequest(v validator.IValidator) {
	validationRules := map[string]string{
		"NewPassword":     "required,min=8",
		"ConfirmPassword": "required,min=8",
		"ResetToken":      "required",
	}
	v.RegisterRules(validationRules, &authapi.ResetPasswordRequest{})
}

func SetRulesOfSignUpRequest(v validator.IValidator) {
	validationRules := map[string]string{
		"Email":     "required,email",
		"Password":  "required,min=8",
		"Phone":     "required,min=9",
		"FirstName": "required",
		"LastName":  "required",
	}
	v.RegisterRules(validationRules, &authapi.SignUpRequest{})
}
