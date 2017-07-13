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
	log := func(s ...interface{}) { fmt.Println(s...) }

	url := ""
	if len(args) > 0 {
		url = args[0]
	}
	log("GET on", url)
	startTime := time.Now()
	resp, err := http.Get(url)

	completedTime := time.Now().Sub(startTime)
	defer log("GET took", completedTime.String())

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
