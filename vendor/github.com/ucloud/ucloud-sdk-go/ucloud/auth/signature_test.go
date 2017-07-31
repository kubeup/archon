package auth

import (
	"testing"
)

func TestGenerateSignature(t *testing.T) {
	var (
		requestParams = "ActionCreateUHostInstanceCPU2ChargeTypeMonthDiskSpace10ImageIdf43736e1-65a5-4bea-ad2e-8a46e18883c2LoginModePasswordMemory2048NameHost01PasswordVUNsb3VkLmNuPublicKeyucloudsomeone@example.com1296235120854146120Quantity1Regioncn-north-01"
		privateKey    = "46f09bb9fab4f12dfc160dae12273d5332b5debe"
		signature     = "64e0fe58642b75db052d50fd7380f79e6a0211bd"
	)

	sig, err := GenerateSignature(requestParams, privateKey)

	if err != nil {
		t.Fatal("GenrateSignature failed")
	}

	if sig != signature {
		t.Fatalf("Expected %s, got %s", signature, sig)
	}

	requestParams = ""
	sig, err = GenerateSignature(requestParams, privateKey)

	if err == nil {
		t.Fatal("GenrateSignature failed")
	}
}
