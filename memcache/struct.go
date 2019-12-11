package memcache

import (
	"sync"
	"time"
)

type Logger interface {
	Printf(fmt string, args ...interface{})
	Debug(s ...interface{})
}

type ValueType struct {
	Value   interface{}
	Expires int64
}

type CacheType struct {
	cache  map[string]*ValueType
	m      sync.RWMutex
	ticker *time.Ticker
	done   chan struct{}
	logger Logger
}
