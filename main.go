package main

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

var (
	// Declare a function that processes a texting task
	exText = func(to string, text string, wrkr string) {
		fmt.Printf("\n\nSending Text Message to: %v\nMessage -> %v\n---\nProcessed by worker: %v\n\n", to, text, wrkr)
	}

	// Declare a function that processes an emailing task
	exEmail = func(to string, body string, wrkr string) {
		fmt.Printf("\n\nSending Email Message to: %v\nMessage -> %v\n---\nProcessed by worker: %v\n\n", to, body, wrkr)
	}
)

// Define a struct type that will describe a job to be processed
type inItem struct {
	Task func(string, string, string)
	To   string
	Text string
}

// Function worker works on a channel (or queue) of
func procJob(jobs <-chan inItem, idWrker string) {
	for job := range jobs {
		// Execute the task referred to by the inbound job
		job.Task(job.To, job.Text, idWrker)

		// Wait for an elapse of time
		time.Sleep(3 * time.Second)
	}
}

func main() {
	loadErr := loadTstData("files/loaddata.toml")
	if loadErr != nil {
		log.Fatal("ERROR: loading the testing data.  Stopping.")
	}

	// Create a channel (queue) that will contain jobs to be executed
	jobs := make(chan inItem)

	// Create a Waitgroup that indicate if workers are still running
	var wg sync.WaitGroup

	// Spawn three workers
	for i := 0; i < 3; i++ {

		// Create a function that serves as a worker and assign to a variable
		workerWrapper := func(fnName string, jobs <-chan inItem) {
			wrkerName := fmt.Sprintf("WRK %v", fnName)
			fmt.Printf("Starting a worker: %v\n", wrkerName)
			defer wg.Done()
			procJob(jobs, wrkerName)
			fmt.Printf("Worker %v is done.\n", wrkerName)
		}

		wg.Add(1)

		// Spin up a worker function
		go workerWrapper(fmt.Sprintf("%p", &workerWrapper), jobs)
	}

	var procItem inItem

	// Create a bunch of jobs
	for i := 0; i < 25; i++ {

		// Create a random integer to generate a distribution of jobs
		myDist := rand.Intn(100)
		switch {
		case myDist <= 70:
			// Tend to create more text messages jobs than email jobs
			tmpItem := inItem{
				Task: exText,
				To:   loadTestData.Numbers[myDist],
				Text: fmt.Sprintf("This is my text to: %v", loadTestData.Numbers[myDist]),
			}
			procItem = tmpItem
		default:
			tmpItem := inItem{
				Task: exEmail,
				To:   loadTestData.Emails[myDist],
				Text: fmt.Sprintf("This is my email to: %v", loadTestData.Emails[myDist]),
			}
			procItem = tmpItem
		}

		// Add a job to the job queue/channel
		jobs <- procItem
	}

	// Close the "jobs" channel (queue)
	close(jobs)
	fmt.Println("Closed jobs!")

	// Wait until the workers complete
	wg.Wait()
}
