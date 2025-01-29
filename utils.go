package configs

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

func Let(ctx context.Context, dst interface{}, src interface{}) error {
	handler := registry.Get(reflect.TypeOf(dst))
	if handler == nil {
		return nil
	}

	return handler.Let(ctx, dst, src)
}

func let[T any](ctx context.Context, dst interface{}, src interface{}) error {
	if g, ok := src.(Getter[T]); ok {
		v, err := g.Get(ctx)
		if err != nil {
			return err
		}
		src = v
	} else {
		if d, ok := dst.(Importer); ok {
			return d.Import(ctx, src)
		}
	}

	if d, ok := dst.(Letter[T]); ok {
		val := reflect.ValueOf(src)
		tp := val.Type()
		if val.Type().ConvertibleTo(tp) {
			var v T
			value := val.Convert(tp)
			reflect.ValueOf(&v).Elem().Set(value)
			return d.Let(ctx, v)
		}
	}

	return nil
}

// Assign assigns values from src to dst.
func Assign(ctx context.Context, dst interface{}, src map[string]interface{}) error {
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
			kind := field.Kind()
			switch kind {
			case reflect.Interface:
				err := Let(ctx, field.Interface(), value)
				if err != nil {
					return err
				}
			case reflect.Struct:
				if val, ok := value.(map[string]interface{}); ok {
					err := Assign(ctx, field.Addr().Interface(), val)
					if err != nil {
						return err
					}
				}
			default:
				val := reflect.ValueOf(value)
				if val.Type().ConvertibleTo(field.Type()) {
					field.Set(val.Convert(field.Type()))
				}
			}
		}
	}
	return nil
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
