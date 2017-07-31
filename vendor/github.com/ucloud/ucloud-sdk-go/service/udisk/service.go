package udisk

import (
	"net/http"

	"github.com/ucloud/ucloud-sdk-go/ucloud"
	"github.com/ucloud/ucloud-sdk-go/ucloud/service"
)

// UDisk api service for UDisk
type UDisk struct {
	*service.Service
}

// New create a uhost service
func New(config *ucloud.Config) *UDisk {

	service := &service.Service{
		Config:      ucloud.DefaultConfig.Merge(config),
		ServiceName: "UHost",
		APIVersion:  ucloud.APIVersion,

		BaseUrl:    ucloud.APIBaseURL,
		HttpClient: &http.Client{},
	}

	return &UDisk{service}

}
