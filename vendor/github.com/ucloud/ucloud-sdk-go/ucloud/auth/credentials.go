package auth

import (
	"errors"
	"os"
)

var (
	ErrPublicKeyIDNotFound = errors.New("UCLOUD_PUBLIC_KEY not found in environment")
	ErrPrivateKeyNotFound  = errors.New("UCLOUD_PRIVATE_KEY not found in environment")
)

type KeyPair struct {
	// UCloud account public key
	PublicKey string

	// UCloud account private key
	PrivateKey string
}

func LoadKeyPairFromEnv() (KeyPair, error) {

	publicKey := os.Getenv("UCLOUD_PUBLIC_KEY")
	privateKey := os.Getenv("UCLOUD_PRIVATE_KEY")

	if publicKey == "" {
		return KeyPair{}, ErrPublicKeyIDNotFound
	}

	if privateKey == "" {
		return KeyPair{}, ErrPrivateKeyNotFound
	}

	return KeyPair{
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}, nil
}
