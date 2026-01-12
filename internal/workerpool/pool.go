package workerpool

import (
	"context"
	"errors"
	"go-admission-service/internal/metrics"
	"sync"
)

var ErrStopped = errors.New("worker pool stopped")

type Job struct {
	Ctx      context.Context
	Do       func(context.Context) Result
	ResultCh chan<- Result
}

type Result struct {
	Value string
	Err   error
}

type Pool struct {
	jobs chan Job

	wg   sync.WaitGroup
	once sync.Once
	done chan struct{}

	metrics *metrics.Counters
}

func New(workerCount int, queueSize int, m *metrics.Counters) *Pool {
	p := &Pool{
		jobs:    make(chan Job, queueSize),
		done:    make(chan struct{}),
		metrics: m,
	}

	p.wg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go p.worker(i)
	}

	return p

}

func (p *Pool) Submit(ctx context.Context, job Job) error {
	select {
	case p.jobs <- job:
		p.metrics.SetQueueDepth(len(p.jobs))
		return nil

	case <-ctx.Done():
		return ctx.Err()

	case <-p.done:
		return ErrStopped
	}

}

func (p *Pool) Stop() {
	p.once.Do(func() {
		close(p.done)
		close(p.jobs)
		p.wg.Wait()
	})
}
