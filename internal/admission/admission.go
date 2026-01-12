package admission

import (
	"context"
	"errors"
	"go-admission-service/internal/metrics"
	"sync"
)

var ErrClosed = errors.New("admission closed")

type Admission struct {
	slots   chan struct{}
	once    sync.Once
	done    chan struct{}
	metrics *metrics.Counters
}

// New creates an admission gate with fixed capacity
func New(capacity int, m *metrics.Counters) *Admission {
	return &Admission{
		slots:   make(chan struct{}, capacity),
		done:    make(chan struct{}),
		metrics: m,
	}
}

// Acquire blocks until a slot is available, ctx is cancelled,
// or admission is closed
// It returns a release function that MUST be called exactly ones
func (a *Admission) Acquire(ctx context.Context) (func(), error) {
	select {
	case a.slots <- struct{}{}:
		a.metrics.IncActiveAdmissions()

		released := false
		release := func() {
			if released {
				return
			}
			released = true
			a.metrics.DecActiveAdmissions()
			<-a.slots
		}
		return release, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-a.done:
		return nil, ErrClosed
	}
}

func (a *Admission) Close() {
	a.once.Do(func() {
		close(a.done)
	})
}
