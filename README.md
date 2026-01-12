# Go Backend Service with Explicit Admission Control

This project demonstrates how to build a production-grade HTTP service in Go
with explicit capacity control, bounded concurrency, and provable invariants.

## Why this exists

Most backend services fail under load due to:
- unbounded goroutines
- hidden queues
- implicit retries
- lack of cancellation

This service makes all limits explicit.

## System Overview

HTTP Request
  → Admission Gate (bounded capacity)
  → Worker Pool (bounded concurrency)
  → Dependency Call (context-aware)
  → Response or Degradation

## Core Invariants

- Active admissions never exceed capacity
- Worker goroutines are bounded
- Requests always terminate
- Queue depth is observable
- System degrades under load instead of collapsing

## Failure Handling

- Overload → 429 Too Many Requests
- Slow dependency → timeout
- Queue saturation → reject
- Cancellation propagates downward

## Observability

The service exposes internal invariants via `/metrics`:
- active_admissions
- active_workers
- in_flight_requests
- queue_depth

These metrics are used to prove the absence of leaks under load.

## How to Run

```bash
go run ./cmd/server
