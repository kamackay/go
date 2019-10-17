package main

import (
	"fmt"
	"github.com/bclicn/color"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"time"
)

const doc = `
Find Usage:
	- find <string (regex)> (defaults to working path, search all files)
	- find <string (regex)> <path> (defaults to search all files)
	- find <string (regex)> <path> <filePattern (regex)>
`

var (
	newLines, _ = regexp.Compile("(\\n|\\r\\n)")
)

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
	findStr := ""
	match := ".*"
	if len(args) > 0 {
		findStr = args[0]
	}
	if findStr == "--help" || findStr == "--h" || findStr == "" {
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

	find := regexp.MustCompile(".*" + findStr + ".*")
	originRegex := regexp.MustCompile(findStr)

	log("Searching for '"+findStr+"' in", dir, "in files matching", match)
	files := getAllFiles(dir)
	sort.Strings(files)
	for _, f := range files {
		filename, exc := filepath.Abs(f)
		if exc == nil {
			r, err := regexp.MatchString(match, filename)
			if r && err == nil {
				totalFiles++
				func() {
					found, bytes := newContains(filename, find)
					if len(found) > 0 {
						log("Found in", color.Green(filename), "("+humanizeBytes(bytes)+")")
						for _, line := range found {
							log("\t", originRegex.ReplaceAllStringFunc(line, color.Blue))
						}
						numFiles++
					}
					bytesRead += bytes
				}()
			}
		}
	}

	log("\n", numFiles, "of", totalFiles, "files contained the pattern")
	log("\t", humanizeBytes(bytesRead), "read")
	completedTime := time.Now().Sub(startTime)
	// show the user how long it took for the server to respond
	defer log("", completedTime.String())
}

func humanizeBytes(bytes int) string {
	unit := 0
	fBytes := float64(bytes)
	for fBytes > 1024 {
		fBytes /= 1024
		unit++
	}
	return strconv.FormatFloat(fBytes, 'f', 4, 64) + " " +
		[]string{"bytes", "KB", "MB", "GB", "TB"}[unit]
}

func newContains(filename string, find *regexp.Regexp) ([]string, int) {
	lines := make([]string, 0)
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return []string{}, 0
	} else {
		for _, line := range newLines.Split(string(data), -1) {
			if subContains(line, find) {
				lines = append(lines, line)
			}
		}
		return lines, len(data)
	}
}

func subContains(content string, find *regexp.Regexp) bool {
	return find.MatchString(content)
}
