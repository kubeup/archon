package main

import (
	"fmt"

	"github.com/ucloud/ucloud-sdk-go/service/uhost"
	"github.com/ucloud/ucloud-sdk-go/ucloud"
	"github.com/ucloud/ucloud-sdk-go/ucloud/auth"
	"github.com/ucloud/ucloud-sdk-go/ucloud/utils"
)

func main() {

	hostsvc := uhost.New(&ucloud.Config{
		Credentials: &auth.KeyPair{
			PublicKey:  "ucloudsomeone@example.com1296235120854146120",
			PrivateKey: "46f09bb9fab4f12dfc160dae12273d5332b5debe",
		},
		Region:    "cn-north-01",
		ProjectID: "",
	})

	describeParams := uhost.DescribeUHostInstanceParams{
		Region: "cn-north-03",
		Limit:  10,
		Offset: 0,
	}

	response, err := hostsvc.DescribeUHostInstance(&describeParams)
	if err != nil {
		fmt.Println(err)
	}
	utils.DumpVal(response)

	//	createUhostParams := uhost.CreateUHostInstanceParams{
	//
	//		Region: "cn-north-03",
	//		ImageId: "uimage-j4fbrn",
	//		LoginMode: "Password",
	//		Password: "UGFzc3dvcmQx",
	//		CPU: 1,
	//		Memory:2048,
	//		Quantity:1,
	//		Quantity:1,
	//		Count: 1,
	//	}
	//
	//
	//	response, err := hostsvc.CreateUHostInstance(&createUhostParams)
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//
	//	utils.DumpVal(response)
	//
	//
	//	// describeimage
	//	imageparams := uhost.DescribeImageParams{
	//		Region: "cn-north-03",
	//	}
	//
	//
	//	imageresp, err := hostsvc.DescribeImage(&imageparams)
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//
	//	utils.DumpVal(imageresp)
}
