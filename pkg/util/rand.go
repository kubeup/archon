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

package util

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

const (
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	alphabets   = "abcdefghijklmnopqrstuvwkyz"
	uuidBytes   = "abcdef1234567890"
	digits      = "0123456789"
)

var (
	ipRE = regexp.MustCompile("^\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}$")
)

func Rand(choices string, n int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = choices[rand.Intn(len(choices))]
	}
	return string(b)
}

func RandNano() string {
	return fmt.Sprintf("%d", time.Now().UTC().UnixNano())
}

func RandString(n int) string {
	return Rand(letterBytes, n)
}

func RandUUID(n int) string {
	return Rand(uuidBytes, n)
}

func RandDigits(n int) string {
	return Rand(digits, n)
}

func RandPassword(n int) (ret string) {
	choices := []string{alphabets, digits, strings.ToUpper(alphabets), digits}
	for ; n > 0; n -= 1 {
		ret = ret + Rand(choices[n%4], 1)
	}

	return
}

func EnvFromMap(d map[string]interface{}) string {
	s := []string{}
	for k, v := range d {
		k = strings.ToUpper(k)
		switch t := v.(type) {
		case string, int, int64, float64:
			// Handle quotes
			s = append(s, fmt.Sprintf("%s=%v", k, t))
		case bool:
			s2 := "0"
			if t {
				s2 = "1"
			}
			s = append(s, fmt.Sprintf("%s=%v", k, s2))
		}
	}

	return strings.Join(s, " ")
}

func EnvMapFromMap(d map[string]interface{}) map[string]interface{} {
	for k, v := range d {
		k2 := strings.ToUpper(k)
		if k2 != k {
			d[k2] = v
			delete(d, k)
		}
	}

	return d
}

func ShortName(n string) string {
	if ipRE.MatchString(n) {
		return strings.Replace(n, ".", "-", -1)
	}

	m := strings.Split(n, ".")
	if len(m) > 0 {
		return m[0]
	}

	return n
}

func OrEmpty(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}
