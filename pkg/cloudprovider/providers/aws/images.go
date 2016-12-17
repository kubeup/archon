package aws

import (
	"encoding/json"
)

type AWSRegion struct {
	PV  string `json:"pv"`
	HVM string `json:"hvm"`
}

var AWSRegions map[string]AWSRegion

const AWSRegionsJSON = `
{"us-east-1": {"pv": "ami-21732036", "hvm": "ami-7e4f1c69"}, "us-west-1": {"pv": "ami-53195233", "hvm": "ami-161a5176"}, "ap-northeast-2": {"pv": "ami-45e2362b", "hvm": "ami-b0e632de"}, "ap-northeast-1": {"pv": "ami-0b8b2d6a", "hvm": "ami-548d2b35"}, "eu-west-1": {"pv": "ami-dd5917ae", "hvm": "ami-925719e1"}, "cn-north-1": {"hvm": "ami-1ce93d71"}, "ap-southeast-1": {"pv": "ami-812e88e2", "hvm": "ami-172e8874"}, "ap-southeast-2": {"pv": "ami-6d3d000e", "hvm": "ami-3b3c0158"}, "us-west-2": {"pv": "ami-61c56101", "hvm": "ami-47f95d27"}, "us-gov-west-1": {"pv": "ami-46db6327", "hvm": "ami-a8de66c9"}, "ap-south-1": {"pv": "ami-9a6a1ef5", "hvm": "ami-986a1ef7"}, "eu-central-1": {"pv": "ami-3a0ff655", "hvm": "ami-b10ff6de"}, "sa-east-1": {"pv": "ami-60019c0c", "hvm": "ami-0d029f61"}}
`

func init() {
	err := json.Unmarshal([]byte(AWSRegionsJSON), &AWSRegions)
	if err != nil {
		panic(err.Error())
	}
}
