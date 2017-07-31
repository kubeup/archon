package auth

import (
	"os"
	"testing"
)

func TestCredentialNotFound(t *testing.T) {
	os.Clearenv()

	_, err := LoadKeyPairFromEnv()
	if err == nil {
		t.Errorf("expected: not found, actual: found")
	}
}
