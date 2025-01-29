package dynConfigs

import (
	"context"
	"fmt"
	"github.com/adverax/configs"
	"reflect"
	"sync"
)

type Integer interface {
	Get(ctx context.Context) (int64, error)
}

type IntegerField struct {
	config Config
	sync.RWMutex
	value int64
}

func (that *IntegerField) Init(c Config) {
	that.config = c
}

func (that *IntegerField) Get(ctx context.Context) (int64, error) {
	that.config.RLock()
	defer that.config.RUnlock()

	return that.Fetch(ctx)
}

func (that *IntegerField) Set(ctx context.Context, value int64) error {
	that.config.Lock()
	defer that.config.Unlock()

	return that.Let(ctx, value)
}

func (that *IntegerField) Fetch(ctx context.Context) (int64, error) {
	that.RLock()
	defer that.RUnlock()

	return that.value, nil
}

func (that *IntegerField) Let(ctx context.Context, value int64) error {
	that.Lock()
	defer that.Unlock()

	that.value = value
	return nil
}

func (that *IntegerField) String() string {
	that.RLock()
	defer that.RUnlock()

	return fmt.Sprintf("%v", that.value)
}

func NewInteger(value int64) *IntegerField {
	return &IntegerField{value: value}
}

type IntegerTypeHandler struct {
	configs.IntegerTypeHandler
}

func (that *IntegerTypeHandler) New(conf Config) interface{} {
	field := NewInteger(0)
	field.Init(conf)
	return field
}

func init() {
	configs.RegisterHandler(reflect.TypeOf((*configs.Integer)(nil)).Elem(), &IntegerTypeHandler{})
}
