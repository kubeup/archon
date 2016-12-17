package util

import (
	"fmt"
	"reflect"
	"strconv"
)

// Convert a string map to a struct with type assertions along the way.
func MapToStruct(in map[string]string, out interface{}, keyPrefix string) (err error) {
	defer func() {
		if err2 := recover(); err2 != nil {
			err = fmt.Errorf("Error: %+v", err2)
		}
	}()

	if in == nil {
		return fmt.Errorf("Can't convert to struct. in is nil")
	}

	ptrType := reflect.TypeOf(out)
	if ptrType.Kind() != reflect.Ptr {
		return fmt.Errorf("Can't convert to struct. out has to be a pointer to a struct")
	}

	v := reflect.ValueOf(out).Elem()
	t := v.Type()
	if t.Kind() != reflect.Struct {
		return fmt.Errorf("Can't convert to struct. out has to be a pointer to a struct")
	}

	for i := 0; i < t.NumField(); i++ {
		fieldType := t.Field(i)
		fieldValue := v.Field(i)

		tagName := fieldType.Tag.Get("k8s")
		if tagName == "" {
			continue
		}
		keyName := keyPrefix + fieldType.Tag.Get("k8s")
		valueString, ok := in[keyName]
		if !ok {
			continue
		}

		var value interface{}
		switch fieldType.Type.Kind() {
		case reflect.Int:
			value, err = strconv.Atoi(valueString)
		case reflect.Bool:
			value, err = strconv.ParseBool(valueString)
		case reflect.Float32:
			value, err = strconv.ParseFloat(valueString, 32)
			value = float32(value.(float64))
		case reflect.Float64:
			value, err = strconv.ParseFloat(valueString, 64)
		case reflect.String:
			value = valueString
		default:
			err = fmt.Errorf("Unsupported field type: %+v", fieldType.Type)
		}
		if err != nil {
			return fmt.Errorf("Unmatch type of field %s: %+v", fieldType.Name, err)
		}
		fieldValue.Set(reflect.ValueOf(value))
	}

	return nil
}

func StructToMap(in interface{}, out map[string]string, prefix string) (err error) {
	defer func() {
		if err2 := recover(); err2 != nil {
			err = fmt.Errorf("Error: %+v", err2)
		}
	}()

	if out == nil {
		return fmt.Errorf("Can't convert to nil map. out is nil")
	}

	t := reflect.TypeOf(in)
	v := reflect.ValueOf(in)
	if t.Kind() == reflect.Ptr {
		v = v.Elem()
		t = v.Type()
	}

	if t.Kind() != reflect.Struct {
		return fmt.Errorf("Can't convert struct. in has to be a struct or *struct")
	}

	for i := 0; i < t.NumField(); i++ {
		fieldType := t.Field(i)
		fieldValue := v.Field(i)

		tagName := fieldType.Tag.Get("k8s")
		if tagName == "" {
			continue
		}

		key := prefix + tagName
		out[key] = fmt.Sprintf("%v", fieldValue.Interface())
	}

	return nil
}
