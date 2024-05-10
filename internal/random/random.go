package random

import (
	"crypto/rand"

	"github.com/google/uuid"
)

func GenerateRandomBytes(n uint32) ([]byte, error) {
	var b = make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, err
}

func GenerateRandomUUID() string {
	return uuid.Must(uuid.NewRandom()).String()
}
