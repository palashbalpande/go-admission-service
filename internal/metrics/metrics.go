package metrics

import "sync/atomic"

type Counters struct {
	ActiveAdmissions int64
	ActiveWorkers    int64
	InFlightRequests int64
	QueueDepth       int64
}

func (c *Counters) IncActiveAdmissions() {
	atomic.AddInt64(&c.ActiveAdmissions, 1)
}

func (c *Counters) DecActiveAdmissions() {
	atomic.AddInt64(&c.ActiveAdmissions, -1)
}

func (c *Counters) IncWorkers() {
	atomic.AddInt64(&c.ActiveWorkers, 1)
}

func (c *Counters) DecWorkers() {
	atomic.AddInt64(&c.ActiveWorkers, -1)
}

func (c *Counters) IncRequests() {
	atomic.AddInt64(&c.InFlightRequests, 1)
}

func (c *Counters) DecRequests() {
	atomic.AddInt64(&c.InFlightRequests, -1)
}

func (c *Counters) SetQueueDepth(n int) {
	atomic.StoreInt64(&c.QueueDepth, int64(n))
}