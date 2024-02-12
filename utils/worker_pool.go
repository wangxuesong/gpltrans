package utils

import "sync"

type WorkerPool struct {
	workers    []*Worker
	taskQueue  chan Task
	numWorkers int
	wg         sync.WaitGroup
}

type Task func()

type Worker struct {
	id      int
	workers *WorkerPool
}

func NewWorkerPool(numWorkers int, taskQueueSize int) *WorkerPool {
	pool := &WorkerPool{
		workers:    make([]*Worker, numWorkers),
		taskQueue:  make(chan Task, taskQueueSize),
		numWorkers: numWorkers,
	}

	for i := 0; i < numWorkers; i++ {
		worker := &Worker{
			id:      i + 1,
			workers: pool,
		}
		pool.workers[i] = worker
		go worker.start()
	}

	return pool
}

func (p *WorkerPool) Submit(task Task) {
	p.taskQueue <- task
}

func (w *Worker) start() {
	for task := range w.workers.taskQueue {
		w.workers.wg.Add(1)
		func() {
			defer func() {
				w.workers.wg.Done()
				// runtime.GC()
			}()
			task()
		}()
	}
}
