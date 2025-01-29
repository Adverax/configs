package dynConfigs

import (
	"context"
	"github.com/adverax/configs"
	"sync"
)

type Config interface {
	Lock()
	Unlock()
	RLock()
	RUnlock()
}

type BaseConfig struct {
	sync.RWMutex
}

type Clonable interface {
	Clone() interface{}
}

type TypeHandler interface {
	configs.TypeHandler
	Get(ctx context.Context, field interface{}) (interface{}, error)
	New(conf Config) interface{}
}
