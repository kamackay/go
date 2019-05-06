package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const doc = `
Find Usage:
	- find <string> (defaults to working path, search all files)
	- find <string> <path> (defaults to search all files)
	- find <string> <path> <filePattern (regex)>
`

func getAllFiles(dir string) []string {
	var files []string
	fileList, _ := ioutil.ReadDir(dir)
	for _, f := range fileList {
		if !f.IsDir() {
			files = append(files, path.Join(dir, f.Name()))
		} else {
			sub := getAllFiles(path.Join(dir, f.Name()))
			files = append(files, sub...)
			/*for _, subFile := range getAllFiles(path.Join(dir, f.Name())) { files = append(files, subFile) }/**/
		}
	}
	return files
}

func main() {
	log := func(s ...interface{}) { fmt.Println(s...) }
	args := os.Args[1:]
	startTime := time.Now()
	dir, _ := os.Getwd()
	find := ""
	match := ".*"
	if len(args) > 0 {
		find = args[0]
	}
	if find == "--help" || find == "--h" || find == "" {
		log(doc)
		return
	}
	if len(args) > 1 {
		dir = args[1]
	}
	if len(args) > 2 {
		match = args[2]
	}
	numFiles, bytesRead := 0, 0

	log("Searching for '"+find+"' in", dir, "in files matching", match)
	for _, f := range getAllFiles(dir) {
		filename, exc := filepath.Abs(f)
		if exc == nil {
			r, err := regexp.MatchString(match, filename)
			if r && err == nil {
				data, err := ioutil.ReadFile(filename)
				if err != nil {
					log(err)
				} else {
					if strings.Contains(string(data), find) {
						log("Found in", filename)
						numFiles++
					}
					bytesRead += len(data)
				}
			}
		}
	}
	unit := 0
	fBytes := float64(bytesRead)
	for fBytes > 1024 {
		fBytes /= 1024
		unit++
	}
	log("\n", numFiles, "files contained the pattern")
	log("\t", strconv.FormatFloat(fBytes, 'f', 4, 64),
		[]string{"bytes", "KB", "MB", "GB", "TB"}[unit], "read")
	completedTime := time.Now().Sub(startTime)
	// show the user how long it took for the server to respond
	defer log("", completedTime.String())
}
