package dynConfigs

import (
	"context"
	"fmt"
	"github.com/adverax/configs"
	"reflect"
	"sync"
)

type Float interface {
	Get(ctx context.Context) (float64, error)
}

type FloatField struct {
	config Config
	sync.RWMutex
	value float64
}

func (that *FloatField) Init(c Config) {
	that.config = c
}

func (that *FloatField) Get(ctx context.Context) (float64, error) {
	that.config.RLock()
	defer that.config.RUnlock()

	return that.Fetch(ctx)
}

func (that *FloatField) Set(ctx context.Context, value float64) error {
	that.config.Lock()
	defer that.config.Unlock()

	return that.Let(ctx, value)
}

func (that *FloatField) Fetch(ctx context.Context) (float64, error) {
	that.RLock()
	defer that.RUnlock()

	return that.value, nil
}

func (that *FloatField) Let(ctx context.Context, value float64) error {
	that.Lock()
	defer that.Unlock()

	that.value = value
	return nil
}

func (that *FloatField) String() string {
	that.RLock()
	defer that.RUnlock()

	return fmt.Sprintf("%v", that.value)
}

func NewFloat(value float64) *FloatField {
	return &FloatField{value: value}
}

type FloatTypeHandler struct {
	configs.FloatTypeHandler
}

func (that *FloatTypeHandler) Get(ctx context.Context, field interface{}) (interface{}, error) {
	if f, ok := field.(Float); ok {
		return f.Get(ctx)
	}

	return nil, nil
}

func (that *FloatTypeHandler) New(conf Config) interface{} {
	field := NewFloat(0)
	field.Init(conf)
	return field
}

func init() {
	configs.RegisterHandler(reflect.TypeOf((*configs.Float)(nil)).Elem(), &FloatTypeHandler{})
}
