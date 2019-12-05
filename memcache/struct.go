package memcache

import (
	"sync"
	"time"
)

type ValueType struct {
	Value   interface{}
	Expires int64
}

type CacheType struct {
	cache  map[string]*ValueType
	m      sync.RWMutex
	ticker *time.Ticker
	done   chan struct{}
}
