package jwt

import (
	"errors"
	"os"
	"time"

	jwt3rdlib "github.com/golang-jwt/jwt/v5"
)

type IJsonWebTokenizer interface {
	GenerateToken(interface{}) (string, error)
	ParseToken(string) (interface{}, error)
}

type CustomJWTTokenizer struct{}

func (*CustomJWTTokenizer) GenerateToken(data interface{}) (string, error) {
	type TempClaims struct {
		data interface{}
		jwt3rdlib.RegisteredClaims
	}

	claim := TempClaims{
		data,
		jwt3rdlib.RegisteredClaims{
			// A usual scenario is to set the expiration time relative to the current time
			ExpiresAt: jwt3rdlib.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt3rdlib.NewNumericDate(time.Now()),
			NotBefore: jwt3rdlib.NewNumericDate(time.Now()),
		},
	}
	token := jwt3rdlib.NewWithClaims(jwt3rdlib.SigningMethodHS256, claim)
	ss, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	return ss, err
}

func (*CustomJWTTokenizer) ParseToken(token string) (interface{}, error) {
	type TempClaims struct {
		data interface{}
		jwt3rdlib.RegisteredClaims
	}
	jwttoken, err := jwt3rdlib.ParseWithClaims(token, &TempClaims{}, func(token *jwt3rdlib.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})
	if err != nil {
		return nil, err
	}
	if claim, ok := jwttoken.Claims.(*TempClaims); ok {
		return claim.data, nil
	}
	return nil, errors.New("error get claim from token")
}
