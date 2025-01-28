package dynConfigs

import (
	"errors"
	"time"
)

type WatcherBuilder struct {
	watcher *Watcher
}

func NewWatcherBuilder() *WatcherBuilder {
	return &WatcherBuilder{
		watcher: &Watcher{
			done:     make(chan struct{}),
			interval: time.Minute,
			composer: NewComposer(NewFieldFactory()),
		},
	}
}

func (that *WatcherBuilder) WithConfig(config Config) *WatcherBuilder {
	that.watcher.config = config
	return that
}

func (that *WatcherBuilder) WithNewConfig(newConfig func() Config) *WatcherBuilder {
	that.watcher.newConfig = newConfig
	return that
}

func (that *WatcherBuilder) WithLoader(loader Loader) *WatcherBuilder {
	that.watcher.loader = loader
	return that
}

func (that *WatcherBuilder) WithInterval(interval time.Duration) *WatcherBuilder {
	that.watcher.interval = interval
	return that
}

func (that *WatcherBuilder) WithFieldFactory(factory FieldFactory) *WatcherBuilder {
	that.watcher.composer = NewComposer(factory)
	return that
}

func (that *WatcherBuilder) WithLogger(logger Logger) *WatcherBuilder {
	that.watcher.logger = logger
	return that
}

func (that *WatcherBuilder) Build() (*Watcher, error) {
	if err := that.checkRequiredFields(); err != nil {
		return nil, err
	}

	return that.watcher, nil
}

func (that *WatcherBuilder) checkRequiredFields() error {
	if that.watcher.config == nil {
		return ErrRequiredFieldConfig
	}

	if that.watcher.newConfig == nil {
		return ErrRequiredFieldNewConfig
	}

	if that.watcher.loader == nil {
		return ErrRequiredFieldLoader
	}

	if that.watcher.composer == nil {
		return ErrRequiredFieldComposer
	}

	return nil
}

var (
	ErrRequiredFieldNewConfig = errors.New("newConfig is required")
	ErrRequiredFieldLoader    = errors.New("loader is required")
	ErrRequiredFieldComposer  = errors.New("composer is required")
	ErrRequiredFieldConfig    = errors.New("config is required")
)
