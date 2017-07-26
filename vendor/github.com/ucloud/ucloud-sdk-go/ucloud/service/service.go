package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"

	"github.com/ucloud/ucloud-sdk-go/ucloud"
	"github.com/ucloud/ucloud-sdk-go/ucloud/uclouderr"
	"github.com/ucloud/ucloud-sdk-go/ucloud/utils"
)

type Service struct {
	Config      *ucloud.Config
	ServiceName string
	APIVersion  string

	BaseUrl    string
	HttpClient *http.Client
}

func (s *Service) setHttpHeader(req *http.Request) {
	if len(s.Config.HTTPHeader) == 0 {
		return
	}

	for key, value := range s.Config.HTTPHeader {
		req.Header.Set(key, value)
	}
}

func (s *Service) DoRequest(action string, params interface{}, response interface{}) error {
	requestURL, err := s.RequestURL(action, params)
	if err != nil {
		return fmt.Errorf("build request url failed, error: %s", err)
	}

	httpReq, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return fmt.Errorf("new request url failed, error: %s", err)
	}
	s.setHttpHeader(httpReq)
	httpResp, err := s.HttpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("do request url failed, error: %s", err)
	}

	defer httpResp.Body.Close()
	body, err := ioutil.ReadAll(httpResp.Body)

	if err != nil {
		return fmt.Errorf("do request url failed, error: %s", err)
	}

	statusCode := httpResp.StatusCode
	if statusCode >= 400 && statusCode <= 599 {

		uerr := uclouderr.UcloudError{}
		err = json.Unmarshal(body, &uerr)
		return &uclouderr.RequestFailed{
			UcloudError: uerr,
			StatusCode:  statusCode,
		}
	}

	if err = json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("unmarshal url failed, error: %s body: %v", err, string(body))
	}

	retCode := reflect.ValueOf(response).Elem().FieldByName("RetCode").Int()
	if retCode != 0 {
		var resp *uclouderr.UcloudError
		if err = json.Unmarshal(body, &resp); err != nil {
			return fmt.Errorf("unmarshal fault message failed, error:%s", err)
		}

		message := reflect.ValueOf(resp).Elem().FieldByName("Message").String()
		return fmt.Errorf("RetCode:%d Message:%s", retCode, message)
	}

	return nil
}

// RequestURL is fully url of api request
func (s *Service) RequestURL(action string, params interface{}) (string, error) {
	if len(s.BaseUrl) == 0 {
		return "", errors.New("baseUrl is not set")
	}

	commonRequest := ucloud.CommonRequest{
		Action:    action,
		PublicKey: s.Config.Credentials.PublicKey,
		ProjectId: s.Config.ProjectID,
	}

	values := url.Values{}
	utils.ConvertParamsToValues(commonRequest, &values)
	utils.ConvertParamsToValues(params, &values)

	url, err := utils.UrlWithSignature(values, s.BaseUrl, s.Config.Credentials.PrivateKey)
	if err != nil {
		return "", fmt.Errorf("convert params error: %s", err)
	}

	return url, nil
}
