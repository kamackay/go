package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"os"
	"time"
	"bytes"
)

func main() {
	args := os.Args[1:]
	// Just to make the call to log shorter
	log := func(s ...interface{}) { fmt.Println(s...) }

	var url string
	if len(args) > 0 {
		url = args[0]
	} else {
		log("Please input a URL to GET")
		return
	}
	var data string
	if len(args) > 1 {
		data = args[1]
	} else {
		data = "{}"
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(data)))
	req.Header.Set("Content-Type", "application/json")

	log("POST on", url, "\t\tData:", data, "\n")
	startTime := time.Now()
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	completedTime := time.Now().Sub(startTime)
	// show the user how long it took for the server to respond
	defer log("\nPOST on", url, "took", completedTime.String())

	if err != nil {
		log("There was an error in the GET")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log("There was an error while reading the response")
	}
	var responseString string
	for _, element := range body {
		responseString += string(element)
	}
	log(responseString)
}
