package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"time"

	"github.com/bclicn/color"
)

const doc = `
Find Usage:
	- find <string (regex)> (defaults to working path, search all files)
	- find <string (regex)> <path> (defaults to search all files)
	- find <string (regex)> <path> <filePattern (regex)>
`

var (
	newLines, _ = regexp.Compile("(\\n|(\\r\\n))")
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

type line struct {
	text string
	num  int
}

func main() {
	log := func(s ...interface{}) { fmt.Println(s...) }
	args := os.Args[1:]
	args, opts := splitOpts(args)
	caseInsensitive := in(opts, "-i")
	startTime := time.Now()
	dir, _ := os.Getwd()
	findStr := ""
	match := ".*"
	if len(args) > 0 {
		findStr = args[0]
	}
	if in(opts, "--help") || in(opts, "--h") {
		log(doc)
		return
	}
	if len(args) > 1 {
		dir = args[1]
	}
	if len(args) > 2 {
		match = args[2]
	}
	// numFiles, bytesRead := 0, 0

	totalFiles := 0

	var logSem = make(chan int, 1)
	var readSem = make(chan int, runtime.NumCPU())

	caseFlag := ternaryString(caseInsensitive, "(?i)", "")

	find := regexp.MustCompile(caseFlag + ".*" + findStr + ".*")
	originRegex := regexp.MustCompile(caseFlag + findStr)

	log("Searching for '"+findStr+"'"+
		ternaryString(caseInsensitive, " (Case Insensitive)", "")+
		" in", dir, "in files matching", match, "using", runtime.NumCPU(), "threads")
	files := getAllFiles(dir)
	sort.Strings(files)
	for _, f := range files {
		filename, exc := filepath.Abs(f)
		if exc != nil {
			panic(exc)
		}
		r, err := regexp.MatchString(match, filename)
		if r && err == nil {
			totalFiles++
			readSem <- 1
			go func() {
				found, bytes := search(filename, find)
				if len(found) > 0 {
					logSem <- 1
					log("Found in", color.Green(filename), "("+humanizeBytes(bytes)+")")
					for _, line := range found {
						log("  Line #"+strconv.Itoa(line.num+1), "\t",
							truncate(originRegex.ReplaceAllStringFunc(line.text, color.Blue), 250))
					}
					<-logSem
				}
				<-readSem
				// bytesRead += bytes
			}()
		}
	}

	// log("\n", numFiles, "of", totalFiles, "files contained the pattern")
	// log("\t", humanizeBytes(bytesRead), "read")
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

func search(filename string, find *regexp.Regexp) ([]*line, int) {
	lines := make([]*line, 0)
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return []*line{}, 0
	} else {
		for x, l := range newLines.Split(string(data), -1) {
			if subContains(l, find) {
				lines = append(lines, &line{text: l, num: x})
			}
		}
		return lines, len(data)
	}
}

func in(list []string, s string) bool {
	for _, i := range list {
		if i == s {
			return true
		}
	}
	return false
}

func subContains(content string, find *regexp.Regexp) bool {
	return find.MatchString(content)
}

func splitOpts(args []string) ([]string, []string) {
	opts := make([]string, 0)
	optMatch := regexp.MustCompile("-.*")
	i := 0
	for i < len(args) {
		l := args[i]
		if optMatch.MatchString(l) {
			args = del(args, i)
			opts = append(opts, l)
		}
		i = i + 1
	}
	return args, opts
}

func del(l []string, i int) []string {
	return append(l[:i], l[i+1:]...)
}

func ternaryString(val bool, v1 string, v2 string) string {
	if val {
		return v1
	} else {
		return v2
	}
}

func truncate(s string, max int) string {
	ending := color.Yellow("...")
	if len(s) < max {
		return s
	} else {
		return s[:max-len(ending)] + ending
	}
}
