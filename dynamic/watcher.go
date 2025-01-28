package dynConfigs

import (
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

type Watcher struct {
	config    Config
	newConfig func() Config
	loader    Loader
	interval  time.Duration
	composer  *Composer
	logger    Logger
	done      chan struct{}
}

func (that *Watcher) Start() {
	go that.Serve()
}

func (that *Watcher) Close() {
	close(that.done)
}

func (that *Watcher) Serve() {
	for {
		select {
		case <-that.done:
			return
		case <-time.After(that.interval):
		}

		that.refresh()
	}
}

func (that *Watcher) refresh() {
	config := that.newConfig()
	that.composer.Init(config)
	err := that.loader.Load(config)
	if err != nil {
		if !errors.Is(err, configs.ErrDistinct) {
			if that.logger != nil {
				that.logger.WithError(err).Error("error load config")
			}
		}
		return
	}

	that.composer.Assign(that.config, config)
}
