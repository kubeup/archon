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
