package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"gopkg.in/yaml.v2"
)

// Job represents a task fetched from the GitHub API
type Job struct {
	ID      int
	Name    string
	Payload string // Placeholder for job details
}

// Workflow represents a parsed YAML workflow
type Workflow struct {
	Name  string   `yaml:"name"`
	Steps []string `yaml:"steps"`
}

// Worker processes jobs
func Worker(id int, jobs <-chan Job, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		fmt.Printf("Worker %d starting job %d: %s\n", id, job.ID, job.Name)

		// Simulate workflow execution
		executeWorkflow(job)

		fmt.Printf("Worker %d finished job %d\n", id, job.ID)
	}
}

// PollGitHubJobs fetches jobs from the GitHub API
func PollGitHubJobs(jobQueue chan<- Job) {
	for {
		fmt.Println("Polling GitHub for new jobs...")

		// Simulate an API response
		resp, err := http.Get("https://api.github.com/repos/yourrepo/actions/jobs")
		if err != nil {
			fmt.Println("Error polling jobs:", err)
			time.Sleep(10 * time.Second)
			continue
		}
		defer resp.Body.Close()

		// Parse the API response (simulated here)
		body, _ := ioutil.ReadAll(resp.Body)
		var jobs []Job
		if err := json.Unmarshal(body, &jobs); err != nil {
			fmt.Println("Error parsing job data:", err)
			time.Sleep(10 * time.Second)
			continue
		}

		// Enqueue jobs
		for _, job := range jobs {
			jobQueue <- job
		}

		time.Sleep(30 * time.Second) // Poll every 30 seconds
	}
}

// ParseWorkflow parses YAML workflow data
func ParseWorkflow(yamlData string) (*Workflow, error) {
	var wf Workflow
	err := yaml.Unmarshal([]byte(yamlData), &wf)
	if err != nil {
		return nil, err
	}
	return &wf, nil
}

// ExecuteWorkflow simulates executing a workflow
func executeWorkflow(job Job) {
	// Simulate loading a workflow file
	yamlData := `
name: Example Workflow
steps:
  - "echo 'Step 1: Clone Repo'"
  - "echo 'Step 2: Build'"
  - "echo 'Step 3: Test'"
`
	workflow, err := ParseWorkflow(yamlData)
	if err != nil {
		fmt.Printf("Error parsing workflow for job %d: %v\n", job.ID, err)
		return
	}

	fmt.Printf("Executing workflow: %s\n", workflow.Name)
	for _, step := range workflow.Steps {
		fmt.Printf("Executing step: %s\n", step)
		time.Sleep(1 * time.Second) // Simulate step execution
	}
}

func main() {
	// Job queue
	jobQueue := make(chan Job, 10)

	// Worker pool
	numWorkers := 3
	var wg sync.WaitGroup
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go Worker(i, jobQueue, &wg)
	}

	// Start job polling
	go PollGitHubJobs(jobQueue)

	// Simulate sending test jobs for demo purposes
	go func() {
		for i := 1; i <= 5; i++ {
			jobQueue <- Job{ID: i, Name: fmt.Sprintf("Test Job %d", i)}
			time.Sleep(5 * time.Second)
		}
	}()

	// Wait for workers to finish
	wg.Wait()
	close(jobQueue)
	fmt.Println("All jobs processed")
}
