package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"os"
	"time"
)

func main() {
	args := os.Args[1:]
	// Just to make the call to log shorter
	log := func(s ...interface{}) { fmt.Println(s...) }

	url := ""
	if len(args) > 0 {
		url = args[0]
	} else {
		log("Please input a URL to GET")
		return
	}
	log("GET on", url, "\n")
	startTime := time.Now()
	resp, err := http.Get(url)

	completedTime := time.Now().Sub(startTime)
	// show the user how long it took for the server to respond
	defer log("\nGET on", url, "took", completedTime.String())

	if err != nil {
		log("There was an error in the GET")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log("There was an error while reading the response")
	}
	responseString := ""
	for _, element := range body {
		responseString += string(element)
	}
	log(responseString)
}
