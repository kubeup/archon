package udisk

import (
	"github.com/ucloud/ucloud-sdk-go/ucloud"
)

// CreateUDiskSnapshot params of create instances
type CreateUDiskSnapshotParams struct {
	ucloud.CommonRequest

	Region     string
	Zone       string
	UDiskId    string
	Name       string
	ChargeType string
	Quantity   int
	Comment    string
}

type CreateUDiskSnapshotResponse struct {
	ucloud.CommonResponse
	SnapshotId []string
}

func (u *UDisk) CreateUDiskSnapshot(params *CreateUDiskSnapshotParams) (*CreateUDiskSnapshotResponse, error) {
	response := &CreateUDiskSnapshotResponse{}
	err := u.DoRequest("CreateUDiskSnapshot", params, response)

	return response, err
}

// AttachUDiskParams params of deleting udisk
type AttachUDiskParams struct {
	ucloud.CommonRequest

	Region  string
	Zone    string
	UHostId string
	UDiskId string
}

type AttachUDiskResponse struct {
	ucloud.CommonResponse
	UHostId string
	UDiskId string
}

func (u *UDisk) AttachUDisk(params *AttachUDiskParams) (*AttachUDiskResponse, error) {
	response := &AttachUDiskResponse{}
	err := u.DoRequest("AttachUDisk", params, response)

	return response, err
}

// DeleteUDisk params of deleting udisk
type DeleteUDiskParams struct {
	ucloud.CommonRequest

	Region  string
	Zone    string
	UDiskId string
}

type DeleteUDiskResponse struct {
	ucloud.CommonResponse
	SnapshotId []string
}

func (u *UDisk) DeleteUDisk(params *DeleteUDiskParams) (*DeleteUDiskResponse, error) {
	response := &DeleteUDiskResponse{}
	err := u.DoRequest("DeleteUDisk", params, response)

	return response, err
}
