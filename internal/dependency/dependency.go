package dependency

import (
	"context"
	"errors"
	"math/rand"
	"time"
)

var ErrTimeout = errors.New("dependency timeout")

type Client struct {
	min time.Duration
	max time.Duration
}

func New(min, max time.Duration) *Client {
	return &Client{
		min: min,
		max: max,
	}
}

func (c *Client) Call(ctx context.Context) (string, error) {
	// simulate variable latency
	delay := c.min + time.Duration(rand.Int63n(int64(c.max-c.min)))

	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-timer.C:
		// 20% failure rate
		if rand.Intn(5) == 0 {
			return "", errors.New("upstream error")
		}
		return "ok", nil
	case <-ctx.Done():
		return "", ErrTimeout
	}
}