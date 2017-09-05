package ulb

import (
	"github.com/ucloud/ucloud-sdk-go/ucloud"
)

// ULB

type CreateULBParams struct {
	ucloud.CommonRequest
	Region    string
	ProjectId string

	ULBName    string
	Tag        string
	Remark     string
	OuterMode  string
	InnerMode  string
	ChargeType string
}

type CreateULBResponse struct {
	ucloud.CommonResponse

	ULBId string
}

func (u *ULB) CreateULB(params *CreateULBParams) (*CreateULBResponse, error) {
	response := &CreateULBResponse{}
	err := u.DoRequest("CreateULB", params, response)

	return response, err
}

type DeleteULBParams struct {
	ucloud.CommonRequest
	Region    string
	ProjectId string

	ULBId string
}

type DeleteULBResponse struct {
	ucloud.CommonResponse
}

func (u *ULB) DeleteULB(params *DeleteULBParams) (*DeleteULBResponse, error) {
	response := &DeleteULBResponse{}
	err := u.DoRequest("DeleteULB", params, response)

	return response, err
}

type ULBBandwidthType int

const (
	ULBBandwidthTypePrivate ULBBandwidthType = 0
	ULBBandwidthTypeShared  ULBBandwidthType = 1
)

type ULBIPSet struct {
	OperatorName string
	EIP          string
	EIPId        string
}

type SSLBindedTargetSet struct {
	VServerId   string
	VServerName string
	ULBId       string
	ULBName     string
}

type ULBSSLSet struct {
	SSLId              string
	SSLName            string
	SSLType            string
	SSLContent         string
	CreateTime         string
	SSLBindedTargetSet []SSLBindedTargetSet
}

type ULBBackendSet struct {
	BackendId    string
	ResourceType string
	ResourceId   string
	ResourceName string
	PrivateIP    string
	Port         int
	Enabled      int
	Status       int
}

type VServerProtocol string

const (
	VServerProtocolHTTP  VServerProtocol = "HTTP"
	VServerProtocolHTTPS VServerProtocol = "HTTPS"
	VServerProtocolTCP   VServerProtocol = "TCP"
	VServerProtocolUDP   VServerProtocol = "UDP"

	PersistentTypeNone         = "None"
	PersistentTypeServerInsert = "ServerInsert"
	PersistentTypeUserDefined  = "UserDefined"
)

type ULBVServerSet struct {
	VServerId      string
	VServerName    string
	Protocol       VServerProtocol
	FrontendPort   int
	Method         string
	PersistentType string
	PersistentInfo string
	ClientTimeout  int
	Status         int
	SSLSet         []ULBSSLSet
	BackendSet     []ULBBackendSet
}

type ULBSet struct {
	ULBId         string
	ULBName       string
	Name          string
	Tag           string
	Remark        string
	PrivateIP     string
	BandwidthType ULBBandwidthType
	Bandwidth     int
	CreateTime    int
	ExpireTime    int
	Resource      []ucloud.Resource
	IPSet         []ULBIPSet
	VServerSet    []ULBVServerSet
	ULBType       string
}

type DescribeULBParams struct {
	ucloud.CommonRequest
	Region    string
	ProjectId string

	Offset int
	Limit  int
	ULBId  string
}

type DescribeULBResponse struct {
	ucloud.CommonResponse

	TotalCount int
	DataSet    []ULBSet
}

func (u *ULB) DescribeULB(params *DescribeULBParams) (*DescribeULBResponse, error) {
	response := &DescribeULBResponse{}
	err := u.DoRequest("DescribeULB", params, response)

	return response, err
}

// VServer

type VServerListenType string

const (
	VServerListenTypeRequestProxy    VServerListenType = "RequestProxy"
	VServerListenTypePacketsTransmit VServerListenType = "PacketsTransmit"
)

type CreateVServerParams struct {
	ucloud.CommonRequest
	Region    string
	ProjectId string

	ULBId          string
	VServerName    string
	ListenType     VServerListenType
	Protocol       VServerProtocol
	FrontendPort   int
	Method         string
	PersistentType string
	PersistentInfo string
	ClientTimeout  int
}

type CreateVServerResponse struct {
	ucloud.CommonResponse

	VServerId string
}

func (u *ULB) CreateVServer(params *CreateVServerParams) (*CreateVServerResponse, error) {
	response := &CreateVServerResponse{}
	err := u.DoRequest("CreateVServer", params, response)

	return response, err
}

type DeleteVServerParams struct {
	ucloud.CommonRequest
	Region    string
	ProjectId string

	ULBId     string
	VServerId string
}

type DeleteVServerResponse struct {
	ucloud.CommonResponse
}

func (u *ULB) DeleteVServer(params *DeleteVServerParams) (*DeleteVServerResponse, error) {
	response := &DeleteVServerResponse{}
	err := u.DoRequest("DeleteVServer", params, response)

	return response, err
}

type DescribeVServerParams struct {
	ucloud.CommonRequest
	Region string

	ULBId     string
	VServerId string
}

type DescribeVServerResponse struct {
	ucloud.CommonResponse

	TotalCount int
	DataSet    []ULBVServerSet
}

func (u *ULB) DescribeVServer(params *DescribeVServerParams) (*DescribeVServerResponse, error) {
	response := &DescribeVServerResponse{}
	err := u.DoRequest("DescribeVServer", params, response)

	return response, err
}

// Backend

type AllocateBackendParams struct {
	ucloud.CommonRequest
	Region    string
	ProjectId string

	ULBId        string
	VServerId    string
	ResourceType string
	ResourceId   string
	Port         int
	Enabled      int
}

type AllocateBackendResponse struct {
	ucloud.CommonResponse

	BackendId string
}

func (u *ULB) AllocateBackend(params *AllocateBackendParams) (*AllocateBackendResponse, error) {
	response := &AllocateBackendResponse{}
	err := u.DoRequest("AllocateBackend", params, response)

	return response, err
}

type ReleaseBackendParams struct {
	ucloud.CommonRequest
	Region    string
	ProjectId string

	ULBId     string
	BackendId string
}

type ReleaseBackendResponse struct {
	ucloud.CommonResponse
}

func (u *ULB) ReleaseBackend(params *ReleaseBackendParams) (*ReleaseBackendResponse, error) {
	response := &ReleaseBackendResponse{}
	err := u.DoRequest("ReleaseBackend", params, response)

	return response, err
}
