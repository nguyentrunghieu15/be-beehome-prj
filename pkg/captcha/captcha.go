package captcha

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
)

// ICaptchaService defines the interface for a captcha service
type ICaptchaService interface {
	Validate(token string) error
}

// RecaptchaService implements CaptchaService using goRecaptchav3
type GGRecaptchaService struct {
	secretKey string
}

// NewRecaptchaService creates a new RecaptchaService instance
func (*GGRecaptchaService) Init() *GGRecaptchaService {
	return &GGRecaptchaService{
		secretKey: os.Getenv("GOOGLE_RECAPTCHA_SECRET_KEY"),
	}
}

// Validate verifies a recaptcha token obtained from the client-side
func (s *GGRecaptchaService) Validate(token string) error {
	resp, err := http.PostForm(os.Getenv("GOOGLE_RECAPTCHA_URI_API"),
		map[string][]string{
			"secret":   {s.secretKey},
			"response": {token},
		})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Parse the response
	// You may need to adjust the parsing logic based on the response structure
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	// Check if captcha verification is successful
	success, ok := result["success"].(bool)
	if !ok || !success {
		return errors.New("verify captcha fail")
	}
	return nil
}
