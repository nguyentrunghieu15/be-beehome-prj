package auth

// Data structures for requests and responses
type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type LoginRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
	Captcha  string `json:"captcha"`
}

type LoginResponse struct {
	AccessToken  string `json:"accessToken"`
	ExpireTime   string `json:"expireTime"`
	RefreshToken string `json:"refreshToken"`
	TokenType    string `json:"tokenType"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

type RefreshTokenResponse struct {
	AccessToken string `json:"accessToken"`
}

type ResetPasswordRequest struct {
	NewPassword     string `json:"newPassword"     validate:"required"`
	ConfirmPassword string `json:"confirmPassword" validate:"required"`
	ResetToken      string `json:"resetToken"      validate:"required"`
}

type SignUpRequest struct {
	Email     string `json:"email"     validate:"required,email"`
	Password  string `json:"password"  validate:"required"`
	Phone     string `json:"phone"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}
