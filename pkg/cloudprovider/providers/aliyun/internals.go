/*
Copyright 2016 The Archon Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package aliyun

import (
	"errors"
	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/ecs"
	"github.com/golang/glog"
	"regexp"
	"strings"
)

var (
	aliyunSafeRE = regexp.MustCompile("AccessKeyId=[^:]*")
)

type LogWriter struct {
}

func (lr *LogWriter) Write(p []byte) (int, error) {
	glog.V(2).Infof("%v", string(p))
	return len(p), nil
}

func isNotFound(err error) bool {
	e, ok := err.(*common.Error)
	if !ok {
		return false
	}

	return e.StatusCode == 404 || strings.Index(strings.ToLower(e.Message), "not found") > -1
}

func firstIP(i ecs.IpAddressSetType) string {
	if len(i.IpAddress) == 0 {
		return ""
	}

	return i.IpAddress[0]
}

func aliyunSafeError(err error) error {
	if err == nil {
		return nil
	}

	e, ok := err.(*common.Error)
	if ok {
		e.Message = aliyunSafeRE.ReplaceAllString(e.Message, "")
		return e
	}

	msg := err.Error()
	newMsg := aliyunSafeRE.ReplaceAllString(msg, "")
	if newMsg != msg {
		return errors.New(newMsg)
	}

	return err
}
