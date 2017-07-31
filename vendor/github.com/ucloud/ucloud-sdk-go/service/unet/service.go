package unet

import (
	"net/http"

	"github.com/ucloud/ucloud-sdk-go/ucloud"
	"github.com/ucloud/ucloud-sdk-go/ucloud/service"
)

type UNet struct {
	*service.Service
}

func New(config *ucloud.Config) *UNet {

	service := &service.Service{
		Config:      ucloud.DefaultConfig.Merge(config),
		ServiceName: "UNet",
		APIVersion:  ucloud.APIVersion,

		BaseUrl:    ucloud.APIBaseURL,
		HttpClient: &http.Client{},
	}

	return &UNet{service}

}
