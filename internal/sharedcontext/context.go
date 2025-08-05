package sharedcontext

import "sync"

type SharedContext struct {
    mu    sync.RWMutex
    store map[string]string
}

func New() *SharedContext {
    return &SharedContext{store: make(map[string]string)}
}

func (sc *SharedContext) Set(key, value string) {
    sc.mu.Lock()
    sc.store[key] = value
    sc.mu.Unlock()
}

func (sc *SharedContext) Get(key string) (string, bool) {
    sc.mu.RLock()
    v, ok := sc.store[key]
    sc.mu.RUnlock()
    return v, ok
}
