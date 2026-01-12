package handler

import (
	"fmt"
	"net/http"
	"sync/atomic"

	"go-admission-service/internal/metrics"
)

func MetricsHandler(m *metrics.Counters) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(
			w,
			"active_admissions %d\nactive_workers %d\nin_flight_requests %d\nqueue_depth %d\n",
			atomic.LoadInt64(&m.ActiveAdmissions),
			atomic.LoadInt64(&m.ActiveWorkers),
			atomic.LoadInt64(&m.InFlightRequests),
			atomic.LoadInt64(&m.QueueDepth),
		)
	}
}
