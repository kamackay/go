package main

import (
	"io/ioutil"
	"fmt"
	"os"
	"path"
	"regexp"
	"time"
)

func in(a byte, list []byte) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func getFiles(dir string) []string {
	var files []string
	fileList, _ := ioutil.ReadDir(dir)
	for _, f := range fileList {
		if !f.IsDir() {
			files = append(files, path.Join(dir, f.Name()))
		} else {
			for _, subFile := range getFiles(path.Join(dir, f.Name())) {
				files = append(files, subFile)
			}
		}
	}
	return files
}

func main() {
	log := func(s ...interface{}) { fmt.Println(s...) }
	args := os.Args[1:]
	startTime := time.Now()
	dir, _ := os.Getwd()
	match := ".*py$"
	if len(args) > 0 {
		dir = args[0]
	}
	if len(args) > 1 {
		match = args[1]
	}

	numFiles := 0

	log("Searching for tabs in", dir)
	for _, f := range getFiles(dir) {
		r, err := regexp.MatchString(match, f)
		if r && err == nil {
			data, err := ioutil.ReadFile(f)
			if err != nil {
				log(err)
			} else {
				if in(byte('\t'), data) {
					log("Tab in", f)
					numFiles++
				}
			}
		}
	}
	log("")
	log(numFiles, "files had tabs")
	completedTime := time.Now().Sub(startTime)
	// show the user how long it took for the server to respond
	defer log("\n", completedTime.String())
}
