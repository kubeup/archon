package auth

import (
	"crypto/sha1"
	"errors"
	"fmt"
)

//GenerateSignature generate signature for request string for ucloud request api.
func GenerateSignature(requestString, privateKey string) (string, error) {

	if len(requestString) < 1 {
		return "", errors.New("Wrong request parameters.")
	}

	if len(privateKey) < 1 {
		return "", errors.New("Wrong private key.")
	}

	h := sha1.New()
	h.Write([]byte(requestString + privateKey))

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
