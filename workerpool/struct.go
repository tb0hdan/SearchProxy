package workerpool

import "sync"

type WorkerPool struct {
	Jobs        chan interface{}
	Results     chan interface{}
	WaitGroup   *sync.WaitGroup
	WaitCh      chan struct{}
	AllItemsCh  chan struct{}
	WorkerCount int
	Function    func(interface{}) interface{}
}
