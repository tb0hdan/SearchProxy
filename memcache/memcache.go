package memcache

import "time"

func (mc *CacheType) Get(key string) (value interface{}, ok bool) {
	mc.m.RLock()
	defer mc.m.RUnlock()
	valueType, ok := mc.cache[key]

	if ok {
		value = valueType.Value
	}

	return
}

func (mc *CacheType) Set(key string, value interface{}) {
	mc.SetEx(key, value, 0)
}

func (mc *CacheType) SetEx(key string, value interface{}, expires int64) {
	mc.m.Lock()
	defer mc.m.Unlock()

	if expires > 0 {
		expires += time.Now().Unix()
	}

	mc.cache[key] = &ValueType{
		Value:   value,
		Expires: expires,
	}
}

func (mc *CacheType) Len() (cacheSize int) {
	cacheSize = len(mc.cache)
	return
}

func (mc *CacheType) Cache() (cache map[string]*ValueType) {
	cache = mc.cache
	return
}

func (mc *CacheType) UnsafeDelete(key string) {
	delete(mc.cache, key)
}

func (mc *CacheType) Delete(key string) {
	mc.m.Lock()
	defer mc.m.Unlock()

	mc.UnsafeDelete(key)
}

func (mc *CacheType) Evictor() {
	for {
		select {
		case <-mc.done:
			return
		case <-mc.ticker.C:
			mc.m.Lock()
			for key, value := range mc.cache {
				if value.Expires == 0 {
					continue
				}

				if value.Expires-time.Now().Unix() <= 0 {
					mc.logger.Printf("Evicting %s\n", key)
					mc.UnsafeDelete(key)
				}
			}
			mc.m.Unlock()
		}
	}
}

func (mc *CacheType) Stop() {
	mc.ticker.Stop()
	mc.done <- struct{}{}

	mc.logger.Debug("Memcache is saying goodbye!")
}

func New(logger Logger) (memCache *CacheType) {
	memCache = &CacheType{cache: make(map[string]*ValueType),
		done:   make(chan struct{}),
		ticker: time.NewTicker(1 * time.Second),
		logger: logger,
	}
	go memCache.Evictor()

	return
}
