package main

import (
	"fmt"
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
	- find (defaults to working path, search all files)
	- find <path> 
	- find <path> <filePattern (regex)> 
`

func getAllFiles(dir string) []string {
	fi, _ := os.Stat(dir)
	if !fi.IsDir() {
		return []string{dir}
	}
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

// Return true if the file ends with one of the given extensions
func checkNotExtension(filename string, extensions []string) bool {
	for _, ext := range extensions {
		if strings.HasSuffix(strings.ToLower(filename), strings.ToLower("."+ext)) {
			return true
		}
	}
	return false
}

func readFile(filename string) (string, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	s := string(buf)
	return s, nil
}

func fixFile(filename string) (bool, int, int) {
	if checkNotExtension(filename, []string{"exe", "jar", "apk", "dex"}) {
		return false, 0, 0
	}

	s, err := readFile(filename)
	lineRegx := regexp.MustCompile(`(\n|\r\n)`)

	if err == nil {
		if strings.Contains(s, "\r\n") {
			fmt.Println(filename)
			fmt.Println("\tHad Windows Newlines")
			lines := lineRegx.Split(s, -1)
			fixedText := []byte(strings.Join(lines, "\n"))
			err := ioutil.WriteFile(filename, fixedText, 0644)
			if err != nil {
				fmt.Println(err)
				return false, len(s), 0
			}
			return true, len(s), len(fixedText)
		} else {
			return false, len(s), 0
		}
	} else {
		fmt.Println(err)
	}
	return false, 0, 0
}

func main() {
	log := func(s ...interface{}) { fmt.Println(s...) }
	args := os.Args[1:]
	startTime := time.Now()
	dir, _ := os.Getwd()
	match := ".*"
	if len(args) > 0 {
		dir = args[0]
	}
	if dir == "--help" || dir == "--h" {
		log(doc)
		return
	}
	if len(args) > 1 {
		match = args[1]
	}

	totalFiles := 0
	fixedFiles := 0
	bytesRead := 0
	bytesWritten := 0
	log("Fixing newlines in", dir, "for files matching", match)
	files := getAllFiles(dir)
	sort.Strings(files)

	for _, f := range files {
		filename, exc := filepath.Abs(f)
		if exc == nil {
			r, err := regexp.MatchString(match, filename)
			if r && err == nil {
				totalFiles++
				fixed, read, written := fixFile(filename)
				if fixed {
					fixedFiles++
				}
				bytesRead += read
				bytesWritten += written

			}
		}
	}

	log("Checked a total of", totalFiles, "files, and corrected", fixedFiles, "files")
	logBytes(bytesRead, "Read")
	logBytes(bytesWritten, "Written")
	defer log("", time.Now().Sub(startTime).String())
}

func logBytes(bytes int, title string) {
	unit := 0
	fBytes := float64(bytes)
	for fBytes > 1024 {
		fBytes /= 1024
		unit++
	}
	fmt.Println("\t", strconv.FormatFloat(fBytes, 'f', 4, 64),
		[]string{"bytes", "KB", "MB", "GB", "TB"}[unit], title)
}
