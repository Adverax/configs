package dynConfigs

import (
	"errors"
	"time"
)

type WatchDogBuilder struct {
	watcher *WatchDog
}

func NewWatchDogBuilder() *WatchDogBuilder {
	return &WatchDogBuilder{
		watcher: &WatchDog{
			done:      make(chan struct{}),
			interval:  time.Minute,
			onUpdated: func(bool) {},
		},
	}
}

func (that *WatchDogBuilder) WithOnUpdated(onUpdated func(bool)) *WatchDogBuilder {
	that.watcher.onUpdated = onUpdated
	return that
}

func (that *WatchDogBuilder) WithConfig(config Config) *WatchDogBuilder {
	that.watcher.config = config
	return that
}

func (that *WatchDogBuilder) WithNewConfig(newConfig func() Config) *WatchDogBuilder {
	that.watcher.newConfig = newConfig
	return that
}

func (that *WatchDogBuilder) WithLoader(loader Loader) *WatchDogBuilder {
	that.watcher.loader = loader
	return that
}

func (that *WatchDogBuilder) WithInterval(interval time.Duration) *WatchDogBuilder {
	that.watcher.interval = interval
	return that
}

func (that *WatchDogBuilder) WithLogger(logger Logger) *WatchDogBuilder {
	that.watcher.logger = logger
	return that
}

func (that *WatchDogBuilder) Build() (*WatchDog, error) {
	if err := that.checkRequiredFields(); err != nil {
		return nil, err
	}

	return that.watcher, nil
}

func (that *WatchDogBuilder) checkRequiredFields() error {
	if that.watcher.config == nil {
		return ErrRequiredFieldConfig
	}

	if that.watcher.newConfig == nil {
		return ErrRequiredFieldNewConfig
	}

	if that.watcher.loader == nil {
		return ErrRequiredFieldLoader
	}

	return nil
}

var (
	ErrRequiredFieldNewConfig = errors.New("newConfig is required")
	ErrRequiredFieldLoader    = errors.New("loader is required")
	ErrRequiredFieldConfig    = errors.New("config is required")
)
