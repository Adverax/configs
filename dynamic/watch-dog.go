package dynConfigs

import (
	"context"
	"errors"
	"github.com/adverax/configs"
	"time"
)

type Loader interface {
	Load(config interface{}) error
}

type Logger interface {
	WithError(err error) Logger
	Error(msg string)
}

type WatchDog struct {
	config    Config
	newConfig func() Config
	loader    Loader
	interval  time.Duration
	logger    Logger
	done      chan struct{}
	onUpdated func(bool)
}

func (that *WatchDog) Start() {
	go that.Serve()
}

func (that *WatchDog) Close() {
	close(that.done)
}

func (that *WatchDog) Serve() {
	for {
		select {
		case <-that.done:
			return
		case <-time.After(that.interval):
		}

		that.refresh(context.Background())
	}
}

func (that *WatchDog) refresh(ctx context.Context) {
	config := Init(that.newConfig())
	err := that.loader.Load(config)
	if err != nil {
		if !errors.Is(err, configs.ErrDistinct) {
			if that.logger != nil {
				that.logger.WithError(err).Error("error load config")
			}
		}
		return
	}

	isStatic, err := isStaticUpdated(ctx, that.config, config)
	if err != nil {
		if that.logger != nil {
			that.logger.WithError(err).Error("error update config")
		}
		return
	}

	Assign(that.config, config)
	that.onUpdated(isStatic)
}
