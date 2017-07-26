package uhost

import (
	"net/http"

	"github.com/ucloud/ucloud-sdk-go/ucloud"
	"github.com/ucloud/ucloud-sdk-go/ucloud/service"
)

// UHost api service for UHost
type UHost struct {
	*service.Service
}

// New create a uhost service
func New(config *ucloud.Config) *UHost {

	service := &service.Service{
		Config:      ucloud.DefaultConfig.Merge(config),
		ServiceName: "UHost",
		APIVersion:  ucloud.APIVersion,

		BaseUrl:    ucloud.APIBaseURL,
		HttpClient: &http.Client{},
	}

	return &UHost{service}

}
