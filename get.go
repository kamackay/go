package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"os"
)

func main() {
	args := os.Args[1:]
	log := func(s string) { fmt.Println(s) }

	url := "http://google.com/"
	if len(args) > 0 {
		url = args[0]
	}
	log("GET on " + url)
	resp, err := http.Get(url)

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
