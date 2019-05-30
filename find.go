package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
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

	totalFiles := 0

	log("Searching for '"+find+"' in", dir, "in files matching", match)
	files := getAllFiles(dir)
	sort.Strings(files)
	for _, f := range files {
		filename, exc := filepath.Abs(f)
		if exc == nil {
			r, err := regexp.MatchString(match, filename)
			if r && err == nil {
				totalFiles++
				func() {
					found, bytes := oldContains(filename, find)
					if found {
						log("Found in", filename)
						numFiles++
					}
					bytesRead += bytes
				}()
			}
		}
	}
	unit := 0
	fBytes := float64(bytesRead)
	for fBytes > 1024 {
		fBytes /= 1024
		unit++
	}
	log("\n", numFiles, "of", totalFiles, "files contained the pattern")
	log("\t", strconv.FormatFloat(fBytes, 'f', 4, 64),
		[]string{"bytes", "KB", "MB", "GB", "TB"}[unit], "read")
	completedTime := time.Now().Sub(startTime)
	// show the user how long it took for the server to respond
	defer log("", completedTime.String())
}

/**
 * Read file into memory and search for the string
 */
func oldContains(filename string, find string) (bool, int) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return false, 0
	} else {
		return subContains(string(data), find), len(data)
	}
}

func newContains(filename string, find string) (bool, int) {
	file, err := os.Open(filename)
	defer func() {
		if file != nil {
			err := file.Close()
			if err != nil {
				panic(err)
			}
		}
	}()
	if err != nil {
		fmt.Println("Could not read file")
		return false, 0
	}

	reader := bufio.NewReader(file)
	var builder strings.Builder
	length := 0
	for {
		line, isPrefix, err := reader.ReadLine()
		builder.WriteString(string(line) + "\n")
		length += len(line)

		//builder.WriteString(string(line))
		if subContains(builder.String(), find) {
			return true, length
		}

		if !isPrefix || err == io.EOF {
			break
		} else if err != nil {
			fmt.Println(err)
		}
	}
	return subContains(builder.String(), find), length
}

/**
 * Actually check a string to see if the substring exists in it
 */
func subContains(content string, find string) bool {
	return strings.Contains(content, find)
}
