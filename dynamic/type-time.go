package dynConfigs

import (
	"context"
	"sync"
	"time"
)

var TimeFormat = time.RFC3339

type Time interface {
	Get(ctx context.Context) (time.Time, error)
}

type TimeEx interface {
	Time
	Init(c Config)
}

type TimeField struct {
	config Config
	sync.RWMutex
	value time.Time
}

func (that *TimeField) Init(c Config) {
	that.config = c
}

func (that *TimeField) Get(ctx context.Context) (time.Time, error) {
	that.config.RLock()
	defer that.config.RUnlock()

	return that.Fetch(ctx)
}

func (that *TimeField) Set(ctx context.Context, value time.Time) error {
	that.config.Lock()
	defer that.config.Unlock()

	return that.Let(ctx, value)
}

func (that *TimeField) Fetch(ctx context.Context) (time.Time, error) {
	that.RLock()
	defer that.RUnlock()

	return that.value, nil
}

func (that *TimeField) Let(ctx context.Context, value time.Time) error {
	that.Lock()
	defer that.Unlock()

	that.value = value
	return nil
}

func (that *TimeField) Import(ctx context.Context, value interface{}) error {
	if s, ok := value.(string); ok {
		val, err := time.Parse(TimeFormat, s)
		if err != nil {
			return err
		}
		return that.Let(ctx, val)
	}

	return nil
}

func (that *TimeField) String() string {
	that.RLock()
	defer that.RUnlock()

	return that.value.Format(TimeFormat)
}

func NewTime(value time.Time) *TimeField {
	return &TimeField{value: value}
}
