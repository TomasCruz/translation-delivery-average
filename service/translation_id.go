package service

import (
	"encoding/hex"
	"math/rand"
)

func NewTranslationID() (string, error) {
	b := make([]byte, 10) //equals 20 characters
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}
