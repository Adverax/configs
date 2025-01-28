package dynConfigs

import (
	"context"
	"fmt"
	"sync"
)

type Boolean interface {
	Get(ctx context.Context) (bool, error)
}

type BooleanEx interface {
	Boolean
	Init(c Config)
}

type BooleanField struct {
	config Config
	sync.RWMutex
	value bool
}

func (that *BooleanField) Init(c Config) {
	that.config = c
}

func (that *BooleanField) Get(ctx context.Context) (bool, error) {
	that.config.RLock()
	defer that.config.RUnlock()

	return that.Fetch(ctx)
}

func (that *BooleanField) Set(ctx context.Context, value bool) error {
	that.config.Lock()
	defer that.config.Unlock()

	return that.Let(ctx, value)
}

func (that *BooleanField) Fetch(ctx context.Context) (bool, error) {
	that.RLock()
	defer that.RUnlock()

	return that.value, nil
}

func (that *BooleanField) Let(ctx context.Context, value bool) error {
	that.Lock()
	defer that.Unlock()

	that.value = value
	return nil
}

func (that *BooleanField) String() string {
	that.RLock()
	defer that.RUnlock()

	return fmt.Sprintf("%v", that.value)
}

func NewBoolean(value bool) *BooleanField {
	return &BooleanField{value: value}
}
