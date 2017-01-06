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
	"github.com/stretchr/testify/suite"
	"testing"
)

type testStruct struct {
	A int     `k8s:"a"`
	B float32 `k8s:"b"`
	C bool    `k8s:"c"`
	D string  `k8s:"d"`
	E float64 `k8s:"e"`
}

type M2SSuite struct {
	suite.Suite
}

func (t *M2SSuite) TestMapToStruct() {
	m := map[string]string{
		"a": "1",
		"b": "3.3",
		"c": "true",
		"d": "goodjob",
	}

	s := testStruct{}
	err := MapToStruct(m, &s, "")
	t.NoError(err)
	t.Equal(s.A, 1)
	t.Equal(s.B, float32(3.3))
	t.Equal(s.C, true)
	t.Equal(s.D, "goodjob")
}

func (t *M2SSuite) TestMStructToMap() {
	n := map[string]string{
		"a": "1",
		"b": "3.3",
		"c": "true",
		"d": "goodjob",
		"e": "6.4",
	}
	m := make(map[string]string)

	s := testStruct{
		A: 1,
		B: 3.3,
		C: true,
		D: "goodjob",
		E: 6.4,
	}
	err := StructToMap(&s, m, "")
	t.NoError(err)
	t.Equal(n, m)

	err = StructToMap(s, m, "")
	t.NoError(err)
	t.Equal(n, m)
}

func TestRunM2SSuite(t *testing.T) {
	suite.Run(t, new(M2SSuite))
}
