package uhost

import (
	"github.com/ucloud/ucloud-sdk-go/ucloud"
)

// CreateUHostInstanceParams params of create instances
type CreateUHostInstanceParams struct {
	ucloud.CommonRequest

	Region          string
	Zone            string
	ImageId         string
	LoginMode       string
	Password        string
	KeyPair         string
	CPU             int
	Memory          int
	DiskSpace       int
	Name            string
	NetworkId       string
	SecurityGroupId string
	ChargeType      string
	Quantity        int
	Count           int
	UHostType       string
	NetCapability   string
	Tag             string
	CouponId        string
	BootDiskSpace   int
}

type CreateUHostInstanceResponse struct {
	ucloud.CommonResponse
	UHostIds []string
}

func (u *UHost) CreateUHostInstance(params *CreateUHostInstanceParams) (*CreateUHostInstanceResponse, error) {
	response := &CreateUHostInstanceResponse{}
	err := u.DoRequest("CreateUHostInstance", params, response)

	return response, err
}

type DescribeUHostInstanceParams struct {
	Region   string
	Zone     string
	UHostIds []string
	Tag      string
	Offset   int
	Limit    int
}

type DiskSet struct {
	Type   string
	DiskId string
	Name   string
	Drive  string
	Size   int
}

type DiskSetArray []DiskSet

type IPSet struct {
	Type      string
	IPId      string
	IP        string
	bandwidth int
}

type IPSetArray []IPSet

type UHostSet struct {
	UHostId            string
	UHostType          string
	Zone               string
	StorageType        string
	ImageId            string
	BasicImageId       string
	BasicImageName     string
	Tag                string
	Remark             string
	Name               string
	State              string
	CreateTime         int
	ChargeType         string
	ExpireTime         int
	CPU                int
	Memory             int
	AutoRenew          string
	DiskSet            DiskSetArray
	IPSet              IPSetArray
	NetCapability      string
	NetworkState       string
	TimemachineFeature string
	HotplugFeature     bool
}

type UHostSetArray []UHostSet

type DescribeUHostInstanceResponse struct {
	ucloud.CommonResponse

	TotalCount int
	UHostSet   UHostSetArray
}

func (u *UHost) DescribeUHostInstance(params *DescribeUHostInstanceParams) (*DescribeUHostInstanceResponse, error) {
	response := &DescribeUHostInstanceResponse{}
	err := u.DoRequest("DescribeUHostInstance", params, response)

	return response, err
}

type TerminateUHostInstanceParams struct {
	ucloud.CommonRequest

	Region  string
	Zone    string
	UHostId string
	Destroy int
}

type TerminateUHostInstanceResponse struct {
	ucloud.CommonResponse
	UhostIds []string
}

func (u *UHost) TerminateUHostInstance(params *TerminateUHostInstanceParams) (*TerminateUHostInstanceResponse, error) {
	response := &TerminateUHostInstanceResponse{}
	err := u.DoRequest("TerminateUHostInstance", params, response)

	return response, err
}

type ResizeUHostInstanceParams struct {
	ucloud.CommonRequest

	Region        string
	Zone          string
	UHostId       string
	CPU           int
	Memory        int
	DiskSpace     int
	BootDiskSpace int
	NetCapValue   int
}

type ResizeUHostInstanceResponse struct {
	ucloud.CommonResponse

	UHostId string
}

func (u *UHost) ResizeUHostInstance(params *ResizeUHostInstanceParams) (*ResizeUHostInstanceResponse, error) {
	response := &ResizeUHostInstanceResponse{}
	err := u.DoRequest("ResizeUHostInstance", params, response)

	return response, err
}

type ReinstallUHostInstanceParams struct {
	ucloud.CommonRequest

	Region  string
	Zone    string
	UHostId string

	Password    string
	ImageId     string
	ReserveDisk string
}

type ReinstallUHostInstanceResponse struct {
	ucloud.CommonResponse

	UHostId string
}

func (u *UHost) ReinstallUHostInstance(params *ReinstallUHostInstanceParams) (*ReinstallUHostInstanceResponse, error) {
	response := &ReinstallUHostInstanceResponse{}
	err := u.DoRequest("ReinstallUHostInstance", params, response)

	return response, err
}

type StartUHostInstanceParams struct {
	ucloud.CommonRequest

	Region  string
	Zone    string
	UHostId string
}

type StartUHostInstanceResponse struct {
	ucloud.CommonResponse
	UhostId string
}

func (u *UHost) StartUHostInstance(params *StartUHostInstanceParams) (*StartUHostInstanceResponse, error) {
	response := &StartUHostInstanceResponse{}
	err := u.DoRequest("StartUHostInstance", params, response)

	return response, err
}

type StopUHostInstanceParams struct {
	ucloud.CommonRequest

	Region  string
	Zone    string
	UHostId string
}

type StopUHostInstanceResponse struct {
	ucloud.CommonResponse

	UhostId string
}

func (u *UHost) StopUHostInstance(params *StopUHostInstanceParams) (*StopUHostInstanceResponse, error) {
	response := &StopUHostInstanceResponse{}
	err := u.DoRequest("StopUHostInstance", params, response)

	return response, err
}

type PoweroffUHostInstanceParams struct {
	ucloud.CommonRequest

	Region  string
	Zone    string
	UHostId string
}

type PoweroffUHostInstanceResponse struct {
	ucloud.CommonResponse

	UhostId string
}

func (u *UHost) PoweroffUHostInstance(params *PoweroffUHostInstanceParams) (*PoweroffUHostInstanceResponse, error) {
	response := &PoweroffUHostInstanceResponse{}
	err := u.DoRequest("PoweroffUHostInstance", params, response)

	return response, err
}

type RebootUHostInstanceParams struct {
	ucloud.CommonRequest

	Region  string
	Zone    string
	UHostId string
}

type RebootUHostInstanceResponse struct {
	ucloud.CommonResponse

	UhostId string
}

func (u *UHost) RebootUHostInstance(params *RebootUHostInstanceParams) (*RebootUHostInstanceResponse, error) {
	response := &RebootUHostInstanceResponse{}
	err := u.DoRequest("RebootUHostInstance", params, response)

	return response, err
}

type ResetUHostInstancePasswordParams struct {
	ucloud.CommonRequest

	Region   string
	Zone     string
	UHostId  string
	Password string
}

type ResetUHostInstancePasswordResponse struct {
	ucloud.CommonResponse

	UhostId string
}

func (u *UHost) ResetUHostInstancePassword(params *ResetUHostInstancePasswordParams) (*ResetUHostInstancePasswordResponse, error) {
	response := &ResetUHostInstancePasswordResponse{}
	err := u.DoRequest("ResetUHostInstancePassword", params, response)

	return response, err
}

type ModifyUHostInstanceNameParams struct {
	ucloud.CommonRequest

	Region  string
	Zone    string
	UHostId string
	Name    string
}

type ModifyUHostInstanceNameResponse struct {
	ucloud.CommonResponse

	UHostId string
}

func (u *UHost) ModifyUHostInstanceName(params *ModifyUHostInstanceNameParams) (*ModifyUHostInstanceNameResponse, error) {
	response := &ModifyUHostInstanceNameResponse{}
	err := u.DoRequest("ModifyUHostInstanceName", params, response)

	return response, err
}

type ModifyUHostInstanceTagParams struct {
	ucloud.CommonRequest

	Region  string
	Zone    string
	UHostId string
	Tag     string
}

type ModifyUHostInstanceTagResponse struct {
	ucloud.CommonResponse

	UHostId string
}

func (u *UHost) ModifyUHostInstanceTag(params *ModifyUHostInstanceTagParams) (*ModifyUHostInstanceTagResponse, error) {
	response := &ModifyUHostInstanceTagResponse{}
	err := u.DoRequest("ModifyUHostInstanceTag", params, response)

	return response, err
}

type ModifyUHostInstanceRemarkParams struct {
	ucloud.CommonRequest

	Region  string
	Zone    string
	UHostId string
	Remark  string
}

type ModifyUHostInstanceRemarkResponse struct {
	ucloud.CommonResponse

	UHostId string
}

func (u *UHost) ModifyUHostInstanceRemark(params *ModifyUHostInstanceRemarkParams) (*ModifyUHostInstanceRemarkResponse, error) {
	response := &ModifyUHostInstanceRemarkResponse{}
	err := u.DoRequest("ModifyUHostInstanceRemark", params, response)

	return response, err
}

type GetUHostInstancePriceParams struct {
	ucloud.CommonRequest

	Region             string
	Zone               string
	ImageId            string
	CPU                int
	Memory             int
	Count              int
	ChargeType         string
	StorageType        string
	DiskSpace          int
	UHostType          string
	NetCapability      string
	TimemachineFeature string
}

type PriceSet struct {
	ChargeType string
	Price      float64
}

type GetUHostInstancePriceResponse struct {
	ucloud.CommonResponse

	UHostId  string
	PriceSet []PriceSet
}

func (u *UHost) GetUHostInstancePrice(params *GetUHostInstancePriceParams) (*GetUHostInstancePriceResponse, error) {
	response := &GetUHostInstancePriceResponse{}
	err := u.DoRequest("GetUHostInstancePrice", params, response)

	return response, err
}

type GetUHostInstanceVncInfoParams struct {
	ucloud.CommonRequest

	Region  string
	Zone    string
	UHostId string
}

type GetUHostInstanceVncInfoResponse struct {
	ucloud.CommonResponse

	UHostId     string
	VncIP       string
	VncPort     int
	VncPassword string
}

func (u *UHost) GetUHostInstanceVncInfo(params *GetUHostInstanceVncInfoParams) (*GetUHostInstanceVncInfoResponse, error) {
	response := &GetUHostInstanceVncInfoResponse{}
	err := u.DoRequest("GetUHostInstanceVncInfo", params, response)

	return response, err
}

type DescribeImageParams struct {
	ucloud.CommonRequest

	Region    string
	Zone      string
	ImageType string
	OsType    string
	ImageId   string
	Offset    int
	Limit     int
}

type ImageSet struct {
	ImageId            string
	ImageName          string
	Zone               string
	OsType             string
	OsName             string
	ImageType          string
	Features           []string
	FuncType           string
	IntegratedSoftware string
	Vendor             string
	Links              string
	State              string
	ImageDescription   string
	CreateTime         int
	ImageSize          int
}

type ImageSetArray []ImageSet

type DescribeImageResponse struct {
	ucloud.CommonResponse

	TotalCount int
	ImageSet   ImageSetArray
}

func (u *UHost) DescribeImage(params *DescribeImageParams) (*DescribeImageResponse, error) {
	response := &DescribeImageResponse{}
	err := u.DoRequest("DescribeImage", params, response)

	return response, err
}

type CreateCustomImageParams struct {
	ucloud.CommonRequest

	Region           string
	Zone             string
	UHostId          string
	ImageName        string
	ImageDescription string
}

type CreateCustomImageResponse struct {
	ucloud.CommonResponse

	ImageId string
}

func (u *UHost) CreateCustomImage(params *CreateCustomImageParams) (*CreateCustomImageResponse, error) {
	response := &CreateCustomImageResponse{}
	err := u.DoRequest("CreateCustomImage", params, response)

	return response, err
}

type TerminateCustomImageParams struct {
	ucloud.CommonRequest

	Region  string
	Zone    string
	ImageId string
}

type TerminateCustomImageResponse struct {
	ucloud.CommonResponse

	ImageId string
}

func (u *UHost) TerminateCustomImage(params *TerminateCustomImageParams) (*TerminateCustomImageResponse, error) {
	response := &TerminateCustomImageResponse{}
	err := u.DoRequest("TerminateCustomImage", params, response)

	return response, err
}

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

func (u *UHost) AttachUDisk(params *AttachUDiskParams) (*AttachUDiskResponse, error) {
	response := &AttachUDiskResponse{}
	err := u.DoRequest("AttachUDisk", params, response)

	return response, err
}

type DetachUDiskParams struct {
	ucloud.CommonRequest

	Region  string
	Zone    string
	UHostId string
	UDiskId string
}

type DetachUDiskResponse struct {
	ucloud.CommonResponse

	UHostId string
	UDiskId string
}

func (u *UHost) DetachUDisk(params *DetachUDiskParams) (*DetachUDiskResponse, error) {
	response := &DetachUDiskResponse{}
	err := u.DoRequest("DetachUDisk", params, response)

	return response, err
}

type CreateUHostInstanceSnapshotParams struct {
	ucloud.CommonRequest

	Region  string
	Zone    string
	UHostId string
}

type CreateUHostInstanceSnapshotResponse struct {
	ucloud.CommonResponse

	UHostId      string
	SnapshotName string
}

func (u *UHost) CreateUHostInstanceSnapshot(params *CreateUHostInstanceSnapshotParams) (*CreateUHostInstanceSnapshotResponse, error) {
	response := &CreateUHostInstanceSnapshotResponse{}
	err := u.DoRequest("CreateUHostInstanceSnapshot", params, response)

	return response, err
}

type DescribeUHostInstanceSnapshotParams struct {
	ucloud.CommonRequest

	Region  string
	Zone    string
	UHostId string
}

type SnapshotSet struct {
	SnapshotName string
	SnapshotTime string
}

type DescribeUHostInstanceSnapshotResponse struct {
	ucloud.CommonResponse

	UHostId     string
	SnapshotSet []SnapshotSet
}

func (u *UHost) DescribeUHostInstanceSnapshot(params *DescribeUHostInstanceSnapshotParams) (*DescribeUHostInstanceSnapshotResponse, error) {
	response := &DescribeUHostInstanceSnapshotResponse{}
	err := u.DoRequest("DescribeUHostInstanceSnapshot", params, response)

	return response, err
}
