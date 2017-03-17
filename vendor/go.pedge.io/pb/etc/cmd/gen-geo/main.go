package main

import (
	"sort"
	"strconv"
	"strings"

	"go.pedge.io/pb/etc/cmd/common"
)

type tmplData struct {
	ElemsAlpha2CodeSorted []*tmplElem
	ElemsAlpha3CodeSorted []*tmplElem
}

type tmplElem struct {
	CountryName        string
	CountryAlpha2Code  string
	CountryAlpha3Code  string
	CountryNumericCode uint32
	CurrencyCode       string
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
		ElemsAlpha2CodeSorted: tmplElemsAlpha2CodeSorted(tmplElems),
		ElemsAlpha3CodeSorted: tmplElemsAlpha3CodeSorted(tmplElems),
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
	countryNumericCode, err := strconv.ParseUint(record[4], 10, 32)
	if err != nil {
		return nil, err
	}
	return &tmplElem{
		CountryName:        record[0],
		CountryAlpha2Code:  strings.ToUpper(record[2]),
		CountryAlpha3Code:  strings.ToUpper(record[3]),
		CountryNumericCode: uint32(countryNumericCode),
		CurrencyCode:       strings.ToUpper(record[14]),
	}, nil
}

func tmplElemsAlpha2CodeSorted(tmplElems []*tmplElem) []*tmplElem {
	s := copyTmplElems(tmplElems)
	sort.Sort(byAlpha2Code(s))
	return s
}

type byAlpha2Code []*tmplElem

func (s byAlpha2Code) Len() int               { return len(s) }
func (s byAlpha2Code) Swap(i int, j int)      { s[i], s[j] = s[j], s[i] }
func (s byAlpha2Code) Less(i int, j int) bool { return s[i].CountryAlpha2Code < s[j].CountryAlpha2Code }

func tmplElemsAlpha3CodeSorted(tmplElems []*tmplElem) []*tmplElem {
	s := copyTmplElems(tmplElems)
	sort.Sort(byAlpha3Code(s))
	return s
}

type byAlpha3Code []*tmplElem

func (s byAlpha3Code) Len() int               { return len(s) }
func (s byAlpha3Code) Swap(i int, j int)      { s[i], s[j] = s[j], s[i] }
func (s byAlpha3Code) Less(i int, j int) bool { return s[i].CountryAlpha3Code < s[j].CountryAlpha3Code }

func copyTmplElems(tmplElems []*tmplElem) []*tmplElem {
	s := make([]*tmplElem, len(tmplElems))
	for i, tmplElem := range tmplElems {
		s[i] = tmplElem
	}
	return s
}
