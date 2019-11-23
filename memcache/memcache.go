package memcache

import (
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)


type ValueType struct {
	Value string
	Expires int64
}

type MemCacheType struct {
	cache map[string]*ValueType
	m     sync.RWMutex
}


func (mc *MemCacheType) Get(key string) (value string, ok bool) {
	mc.m.RLock()
	valueType, ok := mc.cache[key]
	if ok {
		value = valueType.Value
	}
	mc.m.RUnlock()
	return
}


func (mc *MemCacheType) Set(key, value string) {
	mc.SetEx(key, value, 0)
}

func (mc *MemCacheType) SetEx(key, value string, expires int64) {
	mc.m.Lock()
	if expires > 0 {
		expires += time.Now().Unix()
	}
	mc.cache[key] = &ValueType{
		Value:   value,
		Expires: expires,
	}
	mc.m.Unlock()
}


func (mc *MemCacheType) Len() (cache_size int) {
	cache_size = len(mc.cache)
	return
}


func (mc *MemCacheType) Cache() (cache map[string]*ValueType) {
	cache = mc.cache
	return
}


func (mc *MemCacheType) Delete(key string) {
	mc.m.Lock()
	delete(mc.cache, key)
	mc.m.Unlock()
}

func (mc *MemCacheType) Evictor() {
	for {
		for key, value := range mc.cache {
			if value.Expires == 0 {
				continue
			}
			if value.Expires - time.Now().Unix() <= 0 {
				log.Printf("Evicting %s\n", key)
				mc.Delete(key)
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func New() (memCache *MemCacheType) {
	memCache = &MemCacheType{cache: make(map[string]*ValueType)}
	go memCache.Evictor()
	return
}
