package admission

import (
	"context"
	"errors"
	"sync"
)

var ErrClosed = errors.New("admission closed")

type Admission struct {
	slots chan struct{}
	once  sync.Once
	done  chan struct{}
}

// New creates an admission gate with fixed capacity
func New(capacity int) *Admission {
	return &Admission{
		slots: make(chan struct{}, capacity),
		done:  make(chan struct{}),
	}
}

// Acquire blocks until a slot is available, ctx is cancelled,
// or admission is closed
// It returns a release function that MUST be called exactly ones
func (a *Admission) Acquire(ctx context.Context) (func(), error) {
	select {
	case a.slots <- struct{}{}:
		released := false

		release := func() {
			if released {
				return
			}
			released = true
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
