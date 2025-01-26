package configs

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

// Assign assigns values from src to dst.
func Assign(src map[string]interface{}, dst interface{}) {
	dstVal := reflect.ValueOf(dst).Elem()
	dstType := dstVal.Type()

	for i := 0; i < dstVal.NumField(); i++ {
		field := dstVal.Field(i)
		fieldType := dstType.Field(i)
		tag := fieldType.Tag.Get("config")

		if tag == "" {
			tag = strings.ToLower(fieldType.Name)
		}

		if value, ok := src[tag]; ok {
			if field.Kind() == reflect.Struct {
				Assign(value.(map[string]interface{}), field.Addr().Interface())
			} else {
				val := reflect.ValueOf(value)
				if val.Type().ConvertibleTo(field.Type()) {
					field.Set(val.Convert(field.Type()))
				} else {
					field.Set(val)
				}
			}
		}
	}
}

func override(a, b map[string]interface{}) {
	for k, v := range b {
		if av, ok := a[k]; ok {
			if reflect.TypeOf(v) == reflect.TypeOf(av) {
				switch v.(type) {
				case map[string]interface{}:
					override(av.(map[string]interface{}), v.(map[string]interface{}))
				case []interface{}:
					a[k] = v
				default:
					a[k] = v
				}
			} else {
				a[k] = v
			}
		} else {
			a[k] = v
		}
	}
}

func hashOf(data map[string]interface{}) string {
	bs, _ := json.MarshalIndent(data, "", "")
	return digestOf(bs)
}

func digestOf(bs []byte) string {
	return fmt.Sprintf("%x", md5.Sum(bs))
}
