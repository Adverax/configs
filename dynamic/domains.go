package dynConfigs

import "sync"

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
