package dynConfigs

import (
	"context"
	"fmt"
	"github.com/adverax/configs"
	"reflect"
	"sync"
	"time"
)

type Duration interface {
	Get(ctx context.Context) (time.Duration, error)
}

type DurationField struct {
	config Config
	sync.RWMutex
	value time.Duration
}

func (that *DurationField) Init(c Config) {
	that.config = c
}

func (that *DurationField) Get(ctx context.Context) (time.Duration, error) {
	that.config.RLock()
	defer that.config.RUnlock()

	return that.Fetch(ctx)
}

func (that *DurationField) Set(ctx context.Context, value time.Duration) error {
	that.config.Lock()
	defer that.config.Unlock()

	return that.Let(ctx, value)
}

func (that *DurationField) Fetch(ctx context.Context) (time.Duration, error) {
	that.RLock()
	defer that.RUnlock()

	return that.value, nil
}

func (that *DurationField) Let(ctx context.Context, value time.Duration) error {
	that.Lock()
	defer that.Unlock()

	that.value = value
	return nil
}

func (that *DurationField) String() string {
	that.RLock()
	defer that.RUnlock()

	return fmt.Sprintf("%v", that.value)
}

func NewDuration(value time.Duration) *DurationField {
	return &DurationField{value: value}
}

type DurationTypeHandler struct {
	configs.DurationTypeHandler
}

func (that *DurationTypeHandler) Get(ctx context.Context, field interface{}) (interface{}, error) {
	if f, ok := field.(Duration); ok {
		return f.Get(ctx)
	}

	return nil, nil
}

func (that *DurationTypeHandler) New(conf Config) interface{} {
	field := NewDuration(0)
	field.Init(conf)
	return field
}

func init() {
	configs.RegisterHandler(reflect.TypeOf((*configs.Duration)(nil)).Elem(), &DurationTypeHandler{})
}
