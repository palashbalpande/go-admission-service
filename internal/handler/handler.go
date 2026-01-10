package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"go-admission-service/internal/admission"
	"go-admission-service/internal/dependency"
	"go-admission-service/internal/workerpool"
)

type Handler struct {
	admission *admission.Admission
	pool      *workerpool.Pool
	dep       *dependency.Client

	timeout time.Duration
}

func New(
	ad *admission.Admission,
	pool *workerpool.Pool,
	dep *dependency.Client,
	timeout time.Duration,
) *Handler {
	return &Handler{
		admission: ad,
		pool:      pool,
		dep:       dep,
		timeout:   timeout,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), h.timeout)
	defer cancel()

	// if !h.admission.TryAcquire() {
	// 	http.Error(w, "buzy", http.StatusTooManyRequests)
	// 	return
	// }

	// defer h.admission.Release()

	release, err := h.admission.Acquire(ctx)
	if err != nil {
		http.Error(w, "buzy", http.StatusTooManyRequests)
		return
	}

	defer release()

	resultCh := make(chan workerpool.Result, 1)

	job := workerpool.Job{
		Ctx: ctx,
		Do: func(ctx context.Context) workerpool.Result {
			val, err := h.dep.Call(ctx)
			return workerpool.Result{
				Value: val,
				Err:   err,
			}
		},
		ResultCh: resultCh,
	}

	if err := h.pool.Submit(ctx, job); err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	select {
	case res := <-resultCh:
		if res.Err != nil {
			http.Error(w, res.Err.Error(), http.StatusBadGateway)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{
			"result": res.Value,
		})
	case <-ctx.Done():
		http.Error(w, "timeout", http.StatusGatewayTimeout)
	}

}
