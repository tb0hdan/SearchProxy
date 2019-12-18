package workerpool

import "sync"

// WorkerPool - straightforward worker pool implementation with bound methods
type WorkerPool struct {
	Jobs        chan interface{}
	Results     chan interface{}
	WaitGroup   *sync.WaitGroup
	WaitCh      chan struct{}
	AllItemsCh  chan struct{}
	WorkerCount int
	Function    func(interface{}) interface{}
}
