package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	TotalTasks     = 10
	MaxConcurrency = 3 // Simulates 3GB RAM limit
)

func main() {
	//Setting up context for sytem wide graceful cancellations
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//creating a channel to listen to incoming signals from the OS
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	//separate go routine that listens to incoming signals from OS
	go func() {
		sig := <-sigChan
		fmt.Printf("\nReceived signal: %s. Shutting down gracefully...\n", sig)
		cancel() //provided via context
	}()

	// buffered channels
	// Controls how many tasks run at the same time.
	sem := make(chan struct{}, MaxConcurrency)

	var wg sync.WaitGroup

	fmt.Println("--- Service Started  ---")

Loop:
	for i := 1; i <= TotalTasks; i++ {
		select {
		//struct{}{} acts as token (no value but helps in exchange)
		//struct{}is emtpy data object (0 bytes ) , and {} initiliazes it
		//the loop blocks at this place , unless sem has space to acquire a token
		case sem <- struct{}{}:
			// proceeds to go routine
		case <-ctx.Done():
			fmt.Println("Halting new task creation...")
			break Loop
		}

		wg.Add(1)

		//Only happens after lock is acquired)
		go func(taskID int) {
			defer wg.Done()

			defer func() { <-sem }() //release the token when task finished

			processTask(taskID)
		}(i)
	}

	wg.Wait()
	fmt.Println("--- All tasks finished. Exiting. ---")
}

func processTask(id int) {
	fmt.Printf("Task %d has started (Allocating 1GB RAM)\n", id)

	// Simulate work: Random sleep between 5-8 seconds
	sleepTime := time.Duration(rand.Intn(4)+5) * time.Second
	time.Sleep(sleepTime)

	fmt.Printf("Task %d has finished (Releasing RAM)\n", id)
}

