package dynConfigs

import (
	"context"
	"fmt"
	"sync"
)

type Float interface {
	Get(ctx context.Context) (float64, error)
}

type FloatEx interface {
	Float
	Init(c Config)
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
