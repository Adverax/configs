package dynConfigs

import (
	"context"
	"reflect"
	"time"
)

type Composer struct {
	factory FieldFactory
}

func NewComposer(factory FieldFactory) *Composer {
	return &Composer{
		factory: factory,
	}
}

// Init initializes the structure with pointers to the custom types.
func (that *Composer) Init(c Config) {
	that.init(c, c)
}

func (that *Composer) init(c interface{}, conf Config) {
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
				that.newInterface(field, conf)
			} else {
				that.initInterface(field, conf)
			}
		case reflect.Struct:
			that.init(field.Addr().Interface(), conf)
		case reflect.Ptr:
			if field.Elem().Kind() == reflect.Struct {
				that.init(field.Interface(), conf)
				continue
			}
		default:
			continue
		}
	}
}

// Assign copies the values from the source structure to the destination structure.
func (that *Composer) Assign(dst, src Config) {
	dst.Lock()
	defer dst.Unlock()

	src.RLock()
	defer src.RUnlock()

	that.assign(dst, src, dst)
}

func (that *Composer) assign(dst, src interface{}, conf Config) {
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
			_ = that.assignInterface(dstField, srcField, conf)
		case reflect.Struct:
			that.assign(dstField.Addr().Interface(), srcField.Addr().Interface(), conf)
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

			that.assign(dstField.Interface(), srcField.Interface(), conf)
		default:
			dstField.Set(srcField)
		}
	}
}

func (that *Composer) newInterface(field reflect.Value, conf Config) {
	tp := field.Type()
	if tp.Implements(reflect.TypeOf((*Boolean)(nil)).Elem()) {
		val := that.factory.NewBoolean()
		val.Init(conf)
		field.Set(reflect.ValueOf(val))
		return
	}
	if tp.Implements(reflect.TypeOf((*Integer)(nil)).Elem()) {
		val := that.factory.NewInteger()
		val.Init(conf)
		field.Set(reflect.ValueOf(val))
		return
	}
	if tp.Implements(reflect.TypeOf((*Float)(nil)).Elem()) {
		val := that.factory.NewFloat()
		val.Init(conf)
		field.Set(reflect.ValueOf(val))
		return
	}
	if tp.Implements(reflect.TypeOf((*String)(nil)).Elem()) {
		val := that.factory.NewString()
		val.Init(conf)
		field.Set(reflect.ValueOf(val))
		return
	}
	if tp.Implements(reflect.TypeOf((*Duration)(nil)).Elem()) {
		val := that.factory.NewDuration()
		val.Init(conf)
		field.Set(reflect.ValueOf(val))
		return
	}
	if tp.Implements(reflect.TypeOf((*Strings)(nil)).Elem()) {
		val := that.factory.NewStrings()
		val.Init(conf)
		field.Set(reflect.ValueOf(val))
		return
	}
	if tp.Implements(reflect.TypeOf((*Time)(nil)).Elem()) {
		val := that.factory.NewTime()
		val.Init(conf)
		field.Set(reflect.ValueOf(val))
		return
	}
}

func (that *Composer) initInterface(field reflect.Value, conf Config) {
	type Initializer interface {
		Init(conf Config)
	}

	intf := field.Interface()
	if val, ok := intf.(Initializer); ok {
		val.Init(conf)
		return
	}
}

func (that *Composer) assignInterface(dst, src reflect.Value, conf Config) error {
	if dst.IsNil() {
		that.newInterface(dst, conf)
	}

	if dst.Type() != src.Type() {
		return nil
	}

	ctx := context.Background()
	ss := src.Interface()
	dd := dst.Interface()

	switch val := ss.(type) {
	case Boolean:
		return assign[bool](ctx, dd, val)
	case Integer:
		return assign[int64](ctx, dd, val)
	case Float:
		return assign[float64](ctx, dd, val)
	case String:
		return assign[string](ctx, dd, val)
	case Duration:
		return assign[time.Duration](ctx, dd, val)
	case Strings:
		return assign[[]string](ctx, dd, val)
	case Time:
		return assign[time.Time](ctx, dd, val)
	}

	return nil
}

type getter[T any] interface {
	Get(ctx context.Context) (T, error)
}

type letter[T any] interface {
	Let(ctx context.Context, value T) error
}

func assign[T any](ctx context.Context, dst interface{}, src getter[T]) error {
	if vv, ok := dst.(letter[T]); ok {
		v, err := src.Get(ctx)
		if err != nil {
			return err
		}
		return vv.Let(ctx, v)
	}
	return nil
}
