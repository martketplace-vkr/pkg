package workerpool

import (
	"context"
	"sync"

	"github.com/martketplace-vkr/pkg/logger/log"
)

type (
	StateHandler[T any] func(ctx context.Context, data T) error
	Job[T any]          struct {
		Handler StateHandler[T]
		Ctx     context.Context
		Data    T
	}
	StatesPool[T any] struct {
		wg         sync.WaitGroup
		numWorkers int
		jobChan    chan Job[T]
	}

	WorkerPool[T any] interface {
		RunWorkers()
		Shutdown()
		InsertJob(job Job[T])
	}
)

// NewWorkerPool initializes new instance of worker pool.
func NewWorkerPool[T any](numWorkers int) *StatesPool[T] {
	return &StatesPool[T]{
		jobChan:    make(chan Job[T]),
		numWorkers: numWorkers,
	}
}

// worker waits for job from jobChan channel,
// when job is sent to channel, it will be processed.
func (pool *StatesPool[T]) worker() {
	defer pool.wg.Done()

	for job := range pool.jobChan {
		err := job.Handler(job.Ctx, job.Data)
		if err != nil {
			log.Error(err)
		}
	}
}

// RunWorkers runs workers which will be waiting for jobs to handle.
func (pool *StatesPool[T]) RunWorkers() {
	pool.wg.Add(pool.numWorkers)
	for i := 0; i < pool.numWorkers; i++ {
		go pool.worker()
	}
}

// Shutdown closes job channel and
// waits for all workers finished their works.
func (pool *StatesPool[T]) Shutdown() {
	close(pool.jobChan)
	pool.wg.Wait()
}

func (pool *StatesPool[T]) InsertJob(job Job[T]) {
	pool.jobChan <- job
}
