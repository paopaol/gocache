package gocache

import (
	"errors"
	"sync"
)

func newMemoryCache() Cache {
	return &cacheMemory{keyPool: make(map[string][]byte)}
}

type cacheMemory struct {
	mutex   sync.RWMutex
	keyPool map[string][]byte
	Stat
}

func (mem *cacheMemory) Set(key string, value []byte) error {
	mem.mutex.Lock()
	defer mem.mutex.Unlock()

	if v, ok := mem.keyPool[key]; ok {
		mem.delKeyStat(key, v)
	}
	mem.keyPool[key] = value
	mem.addKeyStat(key, value)
	return nil
}

func (mem *cacheMemory) Get(key string) ([]byte, error) {
	mem.mutex.Lock()
	defer mem.mutex.Unlock()

	if v, ok := mem.keyPool[key]; ok {
		return v, nil
	}
	return nil, errors.New("not found key in cache")
}

func (mem *cacheMemory) Del(key string) error {
	mem.mutex.Lock()
	defer mem.mutex.Unlock()

	if v, ok := mem.keyPool[key]; ok {
		delete(mem.keyPool, key)
		mem.delKeyStat(key, v)
	}
	return nil
}

func (mem *cacheMemory) GetStat() Stat {
	mem.mutex.Lock()
	defer mem.mutex.Unlock()

	return mem.Stat
}
