package main

import (
	"fmt"
	"sync"
	"time"
)

// Job represents a task to be executed
type Job struct {
	ID   int
	Name string
}

// Worker processes jobs
func Worker(id int, jobs <-chan Job, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		fmt.Printf("Worker %d starting job %d: %s\n", id, job.ID, job.Name)
		time.Sleep(2 * time.Second) // Simulate job processing
		fmt.Printf("Worker %d finished job %d\n", id, job.ID)
	}
}

func main() {
	jobQueue := make(chan Job, 10) // Buffered channel for job queue
	var wg sync.WaitGroup

	// Start workers
	numWorkers := 3
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go Worker(i, jobQueue, &wg)
	}

	// Send jobs to the queue
	for j := 1; j <= 5; j++ {
		jobQueue <- Job{ID: j, Name: fmt.Sprintf("Job-%d", j)}
	}
	close(jobQueue) // Close channel to signal no more jobs

	// Wait for workers to finish
	wg.Wait()
	fmt.Println("All jobs completed")
}
