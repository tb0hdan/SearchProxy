package workerpool

import "sync"

type WorkerPool struct {
	Jobs        chan interface{}
	Results     chan interface{}
	WaitGroup   *sync.WaitGroup
	WaitCh      chan bool
	AllItemsCh  chan bool
	WorkerCount int
	Function    func(interface{}) interface{}
}
