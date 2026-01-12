

# Strict-Order Task Manager

A concurrent task scheduler in Go that enforces resource limits (max 3 concurrent workers) while guaranteeing strict sequential admission order.

## Architecture

This system uses a **Bounded Semaphore** pattern with **Backpressure** to manage finite resources (simulating a 3GB RAM limit).

### Core Components

1.  **The Producer (Main Loop)**
    * Iterates through tasks `1..N`.
    * **Crucial Logic:** The loop *blocks* at the semaphore acquisition step. It cannot spawn Task `i+1` until it secures a slot for Task `i`. This enforces strict First-Come-First-Served admission.

2.  **The Bouncer (Semaphore)**
    * Implemented as a `chan struct{}` with capacity `3`.
    * Acts as the synchronization barrier between the fast producer and slow consumers.

3.  **The Consumers (Worker Goroutines)**
    * Execute variable-length workloads (simulated 5-8s sleep).
    * Release the semaphore token upon completion, instantly unblocking the Main Loop.

```

## Key points

* **Backpressure:** We block the *generator*, not the *workers*. This prevents "goroutine leaks" where thousands of tasks might sit in memory waiting for a slot.
* **Graceful Shutdown:** Uses `context.Context` to stop accepting new tasks immediately on `SIGINT`, while `sync.WaitGroup` ensures active tasks finish processing before exit.
* **Zero-Allocation Token:** Uses `struct{}{}` for semaphore passing to ensure zero memory overhead for the lock mechanism.

## Usage

```bash
go run main.go

```

To test shutdown:

1. Run the program.
2. Press `Ctrl+C`.
3. Observe that new tasks stop, but active tasks finish safely.



