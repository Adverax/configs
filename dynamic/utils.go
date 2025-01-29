package dynConfigs

import (
	"context"
	"github.com/adverax/configs"
	"reflect"
)

// Init initializes the structure with pointers to the custom types.
func Init(c Config) Config {
	initialize(c, c)
	return c
}

func initialize(c interface{}, conf Config) {
	if reflect.TypeOf(c).Kind() != reflect.Ptr {
		return
	}

	value := reflect.ValueOf(c).Elem()

	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)

		if !field.CanSet() {
			continue
		}

		switch field.Kind() {
		case reflect.Interface:
			if field.IsNil() {
				newInterface(field, conf)
			} else {
				initInterface(field, conf)
			}
		case reflect.Struct:
			initialize(field.Addr().Interface(), conf)
		case reflect.Ptr:
			if field.Elem().Kind() == reflect.Struct {
				initialize(field.Interface(), conf)
				continue
			}
		default:
			continue
		}
	}
}

// Assign copies the values from the source structure to the destination structure.
func Assign(dst, src Config) {
	dst.Lock()
	defer dst.Unlock()

	src.RLock()
	defer src.RUnlock()

	assign(dst, src, dst)
}

func assign(dst, src interface{}, conf Config) {
	dstValue := reflect.ValueOf(dst).Elem()
	srcValue := reflect.ValueOf(src).Elem()

	for i := 0; i < srcValue.NumField(); i++ {
		srcField := srcValue.Field(i)
		dstField := dstValue.Field(i)

		if !srcField.CanInterface() || !dstField.CanSet() {
			continue
		}

		kind := srcField.Kind()
		switch kind {
		case reflect.Interface:
			_ = letInterface(dstField, srcField, conf)
		case reflect.Struct:
			assign(dstField.Addr().Interface(), srcField.Addr().Interface(), conf)
		case reflect.Ptr:
			if srcField.IsNil() {
				continue
			}

			if dstField.IsNil() {
				if c, ok := srcField.Interface().(Clonable); ok {
					dstField.Set(reflect.ValueOf(c.Clone()))
				} else {
					dstField.Set(reflect.New(srcField.Type().Elem()))
				}
			}

			assign(dstField.Interface(), srcField.Interface(), conf)
		default:
			dstField.Set(srcField)
		}
	}
}

func newInterface(field reflect.Value, conf Config) {
	handler := configs.HandlerOf(field.Type())
	if handler != nil {
		if h, ok := handler.(TypeHandler); ok {
			val := h.New(conf)
			field.Set(reflect.ValueOf(val))
		}
	}
}

func initInterface(field reflect.Value, conf Config) {
	type Initializer interface {
		Init(conf Config)
	}

	intf := field.Interface()
	if val, ok := intf.(Initializer); ok {
		val.Init(conf)
		return
	}
}

func letInterface(dst, src reflect.Value, conf Config) error {
	if dst.IsNil() {
		newInterface(dst, conf)
	}

	if dst.Type() != src.Type() {
		return nil
	}

	return configs.Let(context.Background(), dst.Interface(), src.Interface())
}
