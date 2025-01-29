package dynConfigs

import (
	"context"
	"fmt"
	"github.com/adverax/configs"
	"reflect"
	"sync"
)

type Strings interface {
	Get(ctx context.Context) ([]string, error)
}

type StringsField struct {
	config Config
	sync.RWMutex
	value []string
}

func (that *StringsField) Init(c Config) {
	that.config = c
}

func (that *StringsField) Get(ctx context.Context) ([]string, error) {
	that.config.RLock()
	defer that.config.RUnlock()

	return that.Fetch(ctx)
}

func (that *StringsField) Set(ctx context.Context, value []string) error {
	that.config.Lock()
	defer that.config.Unlock()

	return that.Let(ctx, value)
}

func (that *StringsField) Fetch(ctx context.Context) ([]string, error) {
	that.RLock()
	defer that.RUnlock()

	value := make([]string, len(that.value))
	copy(value, that.value)
	return value, nil
}

func (that *StringsField) Let(ctx context.Context, value []string) error {
	that.Lock()
	defer that.Unlock()

	that.value = value
	return nil
}

func (that *StringsField) String() string {
	that.RLock()
	defer that.RUnlock()

	return fmt.Sprintf("%v", that.value)
}

func NewStrings(value []string) *StringsField {
	return &StringsField{value: value}
}

type StringsTypeHandler struct {
	configs.StringsTypeHandler
}

func (that *StringsTypeHandler) New(conf Config) interface{} {
	field := NewStrings(nil)
	field.Init(conf)
	return field
}

func init() {
	configs.RegisterHandler(reflect.TypeOf((*configs.Strings)(nil)).Elem(), &StringsTypeHandler{})
}
