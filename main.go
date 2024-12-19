package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Starting self-hosted GitHub Actions...")

	// Simulate a runner
	for {
		fmt.Println("Polling for jobs...")
		time.Sleep(5 * time.Second)
	}
}
