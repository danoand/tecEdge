package main

import (
	"io/ioutil"
	"log"

	"github.com/BurntSushi/toml"
)

var (
	// Declare a variable to hold test emails and phone numbers
	loadTestData testData
)

// Define a type to models the test data that will be loaded
type testData struct {
	Emails  []string
	Numbers []string
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
