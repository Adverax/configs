package dynConfigs

import (
	"context"
	"fmt"
	"github.com/adverax/configs"
	"reflect"
	"sync"
)

type String interface {
	Get(ctx context.Context) (string, error)
}

type StringField struct {
	config Config
	sync.RWMutex
	value string
}

func (that *StringField) Init(c Config) {
	that.config = c
}

func (that *StringField) Get(ctx context.Context) (string, error) {
	that.config.RLock()
	defer that.config.RUnlock()

	return that.Fetch(ctx)
}

func (that *StringField) Set(ctx context.Context, value string) error {
	that.config.Lock()
	defer that.config.Unlock()

	return that.Let(ctx, value)
}

func (that *StringField) Fetch(ctx context.Context) (string, error) {
	that.RLock()
	defer that.RUnlock()

	return that.value, nil
}

func (that *StringField) Let(ctx context.Context, value string) error {
	that.Lock()
	defer that.Unlock()

	that.value = value
	return nil
}

func (that *StringField) String() string {
	that.RLock()
	defer that.RUnlock()

	return fmt.Sprintf("%v", that.value)
}

func NewString(value string) *StringField {
	return &StringField{value: value}
}

type StringTypeHandler struct {
	configs.StringTypeHandler
}

func (that *StringTypeHandler) New(conf Config) interface{} {
	field := NewString("")
	field.Init(conf)
	return field
}

func init() {
	configs.RegisterHandler(reflect.TypeOf((*configs.String)(nil)).Elem(), &StringTypeHandler{})
}
