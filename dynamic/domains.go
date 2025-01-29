package dynConfigs

import (
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
	New(conf Config) interface{}
}
