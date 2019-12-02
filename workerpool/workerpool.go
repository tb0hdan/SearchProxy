package workerpool

import (
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)


func (wp *WorkerPool) Worker(id int) {
	for j := range wp.Jobs {
		log.Debugf("worker %d started job %v", id, j)
		//
		result := wp.Function(j)
		wp.WaitGroup.Done()
		//
		log.Debugf("worker %d finished job %v", id, j)
		wp.Results <- result

	}
}

func (wp *WorkerPool) ProcessItems(items []interface{}) (results []interface{}) {
	go func() {
		var needBreak bool
		for {
			select {
			case item := <-wp.Results:
				results = append(results, item)
				log.Debug("result!")
			case <-wp.AllItemsCh:
				log.Debug("All items!")
				wp.WaitCh <- true
				needBreak = true
				break
			case <-wp.WaitCh:
				log.Debug("reader got waitch")
				needBreak = true
				break
			}

			if needBreak {
				break
			}

		}
		log.Debug("reader exit")
	}()

	// All items done, signal exit
	go func() {
		wp.WaitGroup.Wait()
		wp.AllItemsCh <- true
		log.Debug("Wooo!")
	}()

	wp.WaitGroup.Add(len(items))
	for w := 1; w <= wp.WorkerCount; w++ {
		go wp.Worker(w)
	}

	for _, item := range items {
		wp.Jobs <- item
	}
	log.Debug("close reached")
	close(wp.Jobs)

	select {
	case <-wp.WaitCh:
		break
	case <-time.After(30 * time.Second):
		break
	}

	log.Debug("return")
	return
}

func New(workerCount int, fn func(item interface{}) interface{}) (wp *WorkerPool) {
	wp = &WorkerPool{}
	wp.Jobs = make(chan interface{})
	wp.Results = make(chan interface{})
	wp.WaitCh = make(chan bool)
	wp.AllItemsCh = make(chan bool)
	wp.WaitGroup = &sync.WaitGroup{}
	wp.WorkerCount = workerCount
	wp.Function = fn
	return wp
}

