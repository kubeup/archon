package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"reflect"
	"sort"
	"strconv"

	"github.com/ucloud/ucloud-sdk-go/ucloud/auth"
)

func ConvertParamsToValues(params interface{}, values *url.Values) {

	elem := reflect.ValueOf(params)
	if elem.Kind() == reflect.Ptr {
		elem = elem.Elem()
	}

	elemType := elem.Type()
	for i := 0; i < elem.NumField(); i++ {
		fieldName := elemType.Field(i).Name

		field := elem.Field(i)
		kind := field.Kind()
		if (kind == reflect.Ptr ||
			kind == reflect.Array ||
			kind == reflect.Slice ||
			kind == reflect.Map ||
			kind == reflect.Chan) && field.IsNil() {
			continue

		}

		if kind == reflect.Ptr {
			field = field.Elem()
			kind = field.Kind()
		}

		var v string
		switch kind {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if field.Int() != 0 {
				v = strconv.FormatInt(field.Int(), 10)
			}

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if field.Uint() != 0 {
				v = strconv.FormatUint(field.Uint(), 10)
			}

		case reflect.Float32:
			v = strconv.FormatFloat(field.Float(), 'f', 4, 32)

		case reflect.Float64:
			v = strconv.FormatFloat(field.Float(), 'f', 4, 64)

		case reflect.Bool:
			v = strconv.FormatBool(field.Bool())

		case reflect.String:
			v = field.String()
		case reflect.Slice:
			switch field.Type().Elem().Kind() {
			case reflect.String:
				l := field.Len()
				if l > 0 {
					for i := 0; i < l; i++ {
						v = field.Index(i).String()
						if v != "" {
							name := elemType.Field(i).Tag.Get("ArgName")
							if name == "" {
								name = fieldName
							}
							name = fmt.Sprintf("%s.%d", name, i)
							values.Set(name, v)
						}
					}
					continue
				}
			default:

			}
		}

		if v != "" {
			name := elemType.Field(i).Tag.Get("ArgName")
			if name == "" {
				name = fieldName
			}

			values.Set(name, v)
		}
	}
}

func UrlWithSignature(values url.Values, baseUrl, privateKey string) (string, error) {

	if values == nil {
		return "", fmt.Errorf("values is empty")
	}

	var buf bytes.Buffer
	keys := make([]string, 0, len(values))
	for k := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		vs := values[k]
		for _, v := range vs {
			buf.WriteString(k)
			buf.WriteString(v)
		}
	}

	signature, err := auth.GenerateSignature(buf.String(), privateKey)
	if err != nil {
		return "", fmt.Errorf("generate signature error:%s", err)
	}

	return baseUrl + "?" + values.Encode() + "&Signature=" + url.QueryEscape(signature), nil
}

func DumpVal(vals ...interface{}) {
	for _, val := range vals {
		prettyJSON, err := json.MarshalIndent(val, "", "    ")
		if err != nil {
			log.Fatal(err)
		}
		log.Print(string(prettyJSON))
	}
}
