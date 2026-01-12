package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-admission-service/internal/admission"
	"go-admission-service/internal/dependency"
	"go-admission-service/internal/handler"
	"go-admission-service/internal/metrics"
	"go-admission-service/internal/workerpool"
)

const (
	addr            = ":8080"
	admissionLimit  = 100
	workerCount     = 10
	queueSize       = 50
	shutdownTimeout = 5 * time.Second
)

func main() {
	log.Println("starting server")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	m := &metrics.Counters{}

	ad := admission.New(admissionLimit, m)

	pool := workerpool.New(workerCount, queueSize, m)

	dep := dependency.New(50*time.Millisecond, 300*time.Millisecond)

	h := handler.New(ad, pool, dep, m, 2*time.Second)

	mux := http.NewServeMux()
	mux.Handle("/", h)
	mux.Handle("/metrics", handler.MetricsHandler(m))

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-sigCh:
		log.Println("shutdown signal received")
	case <-ctx.Done():
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("server shutdown error: %v", err)
	}

	ad.Close()
	pool.Stop()

	log.Println("server stopped cleanly")
}
