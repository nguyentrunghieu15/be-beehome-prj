package random

import (
	"crypto/rand"

	"github.com/google/uuid"
)

func GenerateRandomBytes(n uint32) ([]byte, error) {
	var b = make([]byte, n)
	_, err := rand.Read(b)
	return b, err
}

func GenerateRandomUUID() (string, error) {
	generatedUUID, err := uuid.NewRandom()
	return generatedUUID.String(), err
}
