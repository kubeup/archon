package service

import (
	"net/http"
	"testing"

	"github.com/ucloud/ucloud-sdk-go/ucloud"
	"github.com/ucloud/ucloud-sdk-go/ucloud/auth"
)

type CreateUHostInstance struct {
	Action     string
	Region     string
	ImageId    string
	CPU        int
	Memory     int
	DiskSpace  int
	LoginMode  string
	Password   string
	Name       string
	ChargeType string
	Quantity   int
	PublicKey  string
}

var (
	config = ucloud.Config{
		Credentials: &auth.KeyPair{
			PublicKey:  "ucloudsomeone@example.com1296235120854146120",
			PrivateKey: "46f09bb9fab4f12dfc160dae12273d5332b5debe",
		},
	}

	service = &Service{
		Config:      ucloud.DefaultConfig.Merge(&config),
		ServiceName: "UHost",
		APIVersion:  ucloud.APIVersion,

		BaseUrl:    ucloud.APIBaseURL,
		HttpClient: &http.Client{},
	}

	instance = CreateUHostInstance{
		Action:     "CreateUHostInstance",
		Region:     "cn-north-01",
		ImageId:    "f43736e1-65a5-4bea-ad2e-8a46e18883c2",
		CPU:        2,
		Memory:     2048,
		DiskSpace:  10,
		LoginMode:  "Password",
		Password:   "VUNsb3VkLmNu",
		Name:       "Host01",
		ChargeType: "Month",
		Quantity:   1,
		PublicKey:  "ucloudsomeone@example.com1296235120854146120",
	}
)

const (
	Signature               = "64e0fe58642b75db052d50fd7380f79e6a0211bd"
	RequestUrlWithSignature = "https://api.ucloud.cn/?Action=CreateUHostInstance&CPU=2&ChargeType=Month&DiskSpace=10&ImageId=f43736e1-65a5-4bea-ad2e-8a46e18883c2&LoginMode=Password&Memory=2048&Name=Host01&Password=VUNsb3VkLmNu&PublicKey=ucloudsomeone%40example.com1296235120854146120&Quantity=1&Region=cn-north-01&Signature=64e0fe58642b75db052d50fd7380f79e6a0211bd"
	StringToBeSigned        = "ActionCreateUHostInstanceCPU2ChargeTypeMonthDiskSpace10ImageIdf43736e1-65a5-4bea-ad2e-8a46e18883c2LoginModePasswordMemory2048NameHost01PasswordVUNsb3VkLmNuPublicKeyucloudsomeone@example.com1296235120854146120Quantity1Regioncn-north-0146f09bb9fab4f12dfc160dae12273d5332b5debe"
)

func TestRequestURL(t *testing.T) {

	url, err := service.RequestURL("CreateUHostInstance", instance)
	if err != nil {
		t.Errorf("generate request url failed, Error:%s", err)
	}

	if url != RequestUrlWithSignature {
		t.Errorf("generate url failed. Expect: %s, actual: %s", RequestUrlWithSignature, url)
	}
}
