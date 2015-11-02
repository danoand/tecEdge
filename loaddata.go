package main

import (
	"io/ioutil"
	"log"
	"sync"

	"github.com/BurntSushi/toml"
)

var (
	// Declare a variable to hold test emails and phone numbers
	loadTestData testData

	wrkerMetrics = make(map[string]jobData)
	mutMetrics   = &sync.Mutex{}
	wgMetrics    sync.WaitGroup
)

// Define a type to models the test data that will be loaded
type testData struct {
	Emails  []string
	Numbers []string
}

// jobData defines a struct to hold the number of text and email tasks
//   performed by a worker
type jobData struct {
	NumTexts  int
	NumEmails int
}

// jobIncText increments the number of texts sent for the specified worker
func jobIncText(inKey string) {
	var tNumText int
	var tNumEmails int
	var tJobData jobData

	// Lock wrkerMetrics
	mutMetrics.Lock()

	// Fetch the metrics data
	tNumText = wrkerMetrics[inKey].NumTexts + 1
	tNumEmails = wrkerMetrics[inKey].NumEmails

	// Update new metrics value
	tJobData.NumTexts = tNumText
	tJobData.NumEmails = tNumEmails

	// Assign updated metrics data to map
	wrkerMetrics[inKey] = tJobData

	// Release lock on wrkerMetrics
	mutMetrics.Unlock()

	wgMetrics.Done()
	return
}

// jobIncEmail increments the number of emails sent for the specified worker
func jobIncEmail(inKey string) {
	var tNumText int
	var tNumEmails int
	var tJobData jobData

	// Lock wrkerMetrics
	mutMetrics.Lock()

	// Fetch the metrics data
	tNumText = wrkerMetrics[inKey].NumTexts
	tNumEmails = wrkerMetrics[inKey].NumEmails + 1

	// Update new metrics value
	tJobData.NumTexts = tNumText
	tJobData.NumEmails = tNumEmails

	// Assign updated metrics data to map
	wrkerMetrics[inKey] = tJobData

	// Release lock on wrkerMetrics
	mutMetrics.Unlock()

	wgMetrics.Done()
	return
}

// Function loadRequests loads the requests data structures
func loadTstData(inFName string) (retErr error) {
	var fbytes []byte

	// Read the YAML file data into a byte array
	fbytes, retErr = ioutil.ReadFile(inFName)
	if retErr != nil {
		log.Fatalln("Fn loadTstData: Error reading data to initialize Customers see:", retErr)
	}

	// Unmarshall the YAML data into the data structure
	retErr = toml.Unmarshal(fbytes, &loadTestData)
	if retErr != nil {
		log.Fatalln("Fn loadTstData: Error parsing the Customer TOML data see:", retErr)
	}

	return
}
