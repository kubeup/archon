package main

import (
	"sort"
	"strconv"
	"strings"

	"go.pedge.io/pb/etc/cmd/common"
)

type tmplData struct {
	ElemsCodeSorted []*tmplElem
}

type tmplElem struct {
	CurrencyName        string
	CurrencyCode        string
	CurrencyNumericCode uint32
	CurrencyMinorUnit   uint32
}

func main() {
	common.CSVMain(&generateHelper{})
}

type generateHelper struct{}

func (g *generateHelper) TmplData(records [][]string) (interface{}, error) {
	tmplElems := make([]*tmplElem, len(records)-1)
	for i, record := range records {
		if i == 0 {
			continue
		}
		tmplElem, err := recordToTmplElem(record)
		if err != nil {
			return nil, err
		}
		tmplElems[i-1] = tmplElem
	}
	return &tmplData{
		ElemsCodeSorted: tmplElemsCodeSorted(tmplElems),
	}, nil
}

func (g *generateHelper) ExtraTmplFuncs() map[string]interface{} {
	return map[string]interface{}{
		"cleanCurrencyCode": func(currencyCodeString string) string {
			if currencyCodeString == "" {
				return "_CURRENCY_CODE_NONE"
			}
			return currencyCodeString
		},
	}
}

func recordToTmplElem(record []string) (*tmplElem, error) {
	var currencyMinorUnit uint64
	var currencyNumericCode uint64
	var err error
	if record[16] != "" {
		currencyMinorUnit, err = strconv.ParseUint(record[16], 10, 32)
		if err != nil {
			return nil, err
		}
	}
	if record[18] != "" {
		currencyNumericCode, err = strconv.ParseUint(record[18], 10, 32)
		if err != nil {
			return nil, err
		}
	}
	return &tmplElem{
		CurrencyName:        record[17],
		CurrencyCode:        strings.ToUpper(record[14]),
		CurrencyMinorUnit:   uint32(currencyMinorUnit),
		CurrencyNumericCode: uint32(currencyNumericCode),
	}, nil
}

func tmplElemsCodeSorted(tmplElems []*tmplElem) []*tmplElem {
	filtered := make(map[string]*tmplElem)
	for _, tmplElem := range tmplElems {
		if tmplElem.CurrencyCode != "" {
			filtered[tmplElem.CurrencyCode] = tmplElem
		}
	}
	var s []*tmplElem
	for _, value := range filtered {
		s = append(s, value)
	}
	sort.Sort(byCurrencyCode(s))
	return s
}

type byCurrencyCode []*tmplElem

func (s byCurrencyCode) Len() int               { return len(s) }
func (s byCurrencyCode) Swap(i int, j int)      { s[i], s[j] = s[j], s[i] }
func (s byCurrencyCode) Less(i int, j int) bool { return s[i].CurrencyCode < s[j].CurrencyCode }
