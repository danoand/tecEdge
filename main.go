package main

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

var (
	// Declare a function that processes a "texting" task
	exText = func(to string, text string, wrkr string) {
		// Execute the task
		fmt.Printf("\n\nSending Text Message to: %v\nMessage -> %v\n---\nProcessed by worker: %v\n\n", to, text, wrkr)

		// Increment metrics for the specified worker
		wgMetrics.Add(1)
		go jobIncText(wrkr)
	}

	// Declare a function that processes an "emailing" task
	exEmail = func(to string, body string, wrkr string) {
		// Execute the task
		fmt.Printf("\n\nSending Email Message to: %v\nMessage -> %v\n---\nProcessed by worker: %v\n\n", to, body, wrkr)

		// Increment metrics for the specified worker
		wgMetrics.Add(1)
		go jobIncEmail(wrkr)
	}
)

// inJob defines a struct type that will describe a job to be processed
type inJob struct {
	Task func(string, string, string) // Task to be executed in this job
	To   string                       // The recipient of the outbound message (phone number or email address)
	Text string                       // The outbound message intended for the recipient (text or email)
}

// procjobQueue is a function that "runs" on each worker processing jobQueue
func procJob(jobQueue <-chan inJob, idWrker string) {
	for job := range jobQueue {
		// Execute the task referred to by the inbound job
		job.Task(job.To, job.Text, idWrker)

		// Wait for an elapse of time
		time.Sleep(1 * time.Second)
	}
}

func main() {
	// Load testing data
	loadErr := loadTstData("files/loaddata.toml")
	if loadErr != nil {
		log.Fatal("ERROR: loading the testing data.  Stopping.")
	}

	// Create a channel (queue) that will contain jobQueue to be executed
	jobQueue := make(chan inJob)

	// Create a Waitgroup that indicates if workers are still running
	var wg sync.WaitGroup

	// Spawn three workers
	for i := 0; i < 3; i++ {

		// Create a function that serves as a worker and assign to a variable
		workerWrapper := func(fnName string, jobQueue <-chan inJob) {
			wrkerName := fmt.Sprintf("%v", fnName)
			fmt.Printf("WORKER: Starting a worker: %v\n", wrkerName)
			defer wg.Done()
			procJob(jobQueue, wrkerName)
			fmt.Printf("WORKER: Worker %v is done.\n", wrkerName)
		}

		wg.Add(1)

		// Spin up a worker function
		go workerWrapper(fmt.Sprintf("%p", &workerWrapper), jobQueue)
		// Add an entry to wrkerMetrics to gather metrics on the worker
		wrkerMetrics[fmt.Sprintf("%p", &workerWrapper)] = jobData{0, 0}
	}

	var procJob inJob

	// Create a bunch of jobQueue
	for i := 0; i < 25; i++ {

		// Create a random integer to generate a distribution of jobQueue
		myDist := rand.Intn(100)
		switch {
		case myDist <= 70:
			// Tend to create more text messages jobQueue than email jobQueue
			tmpItem := inJob{
				Task: exText,
				To:   loadTestData.Numbers[myDist],
				Text: fmt.Sprintf("This is my text to: %v", loadTestData.Numbers[myDist]),
			}
			procJob = tmpItem
		default:
			tmpItem := inJob{
				Task: exEmail,
				To:   loadTestData.Emails[myDist],
				Text: fmt.Sprintf("This is my email to: %v", loadTestData.Emails[myDist]),
			}
			procJob = tmpItem
		}

		// Add a job to the job queue/channel
		jobQueue <- procJob
	}

	// Close the "jobQueue" channel (queue)
	close(jobQueue)
	fmt.Println("Closed jobQueue!")

	// Wait until the workers complete
	wg.Wait()

	// Print out metrics to the log
	fmt.Printf("\n\n*********** PRINT OUT METRICS FOR EACH WORKER *************\n")
	// Range through the metrics data
	for k, v := range wrkerMetrics {
		fmt.Printf("WORKER %v processed %v texts & %v emails\n", k, v.NumTexts, v.NumEmails)
	}
	fmt.Printf("************************\n\n")
}
