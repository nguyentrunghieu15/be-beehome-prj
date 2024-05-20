package jwt

import (
	"errors"
	"os"
	"time"

	jwt3rdlib "github.com/golang-jwt/jwt/v5"
)

type IJsonWebTokenizer interface {
	GenerateToken(string, TokenConfigure) (string, error)
	ParseToken(string) (interface{}, error)
}

type TokenConfigure struct {
	ExpiresTime time.Duration
	IssuedAt    time.Time
	NotBefore   time.Time
}

var DefaultAccessTokenConfigure = TokenConfigure{
	ExpiresTime: 3 * time.Hour,
	IssuedAt:    time.Now(),
	NotBefore:   time.Now(),
}

var DefaultRefreshTokenConfigure = TokenConfigure{
	ExpiresTime: 24 * time.Hour,
	IssuedAt:    time.Now(),
	NotBefore:   time.Now(),
}

type CustomJWTTokenizer struct{}

func (*CustomJWTTokenizer) Init() interface{} {
	return &CustomJWTTokenizer{}
}

type TempClaims struct {
	UserId string
	jwt3rdlib.RegisteredClaims
}

func (*CustomJWTTokenizer) GenerateToken(userId string, config TokenConfigure) (string, error) {

	claim := &TempClaims{
		userId,
		jwt3rdlib.RegisteredClaims{
			// A usual scenario is to set the expiration time relative to the current time
			ExpiresAt: jwt3rdlib.NewNumericDate(time.Now().Add(config.ExpiresTime)),
			IssuedAt:  jwt3rdlib.NewNumericDate(config.IssuedAt),
			NotBefore: jwt3rdlib.NewNumericDate(config.NotBefore),
		},
	}

	token := jwt3rdlib.NewWithClaims(jwt3rdlib.SigningMethodHS256, claim)
	ss, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	return ss, err
}

func (*CustomJWTTokenizer) ParseToken(token string) (interface{}, error) {
	jwttoken, err := jwt3rdlib.ParseWithClaims(token, &TempClaims{}, func(token *jwt3rdlib.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})
	if err != nil {
		return nil, err
	}
	if claim, ok := jwttoken.Claims.(*TempClaims); ok {
		return claim.UserId, nil
	}
	return nil, errors.New("error get claim from token")
}
