package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

type ComputedFile struct {
	info os.FileInfo
	size int64
}

func contains(arr []string, s string) bool {
	cleanString := func(st string) string {
		return strings.ToLower(strings.Trim(st, " "))
	}
	lowerS := cleanString(s)
	for _, a := range arr {
		if lowerS == cleanString(a) {
			return true
		}
	}
	return false
}

func dirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}

func repeatChar(c string, n int) string {
	arr := make([]string, n)
	return strings.Join(arr, c)
}

func prettySize(size int64) string {
	var fSize = float64(size)
	var unit uint8 = 0
	for fSize > 1024 {
		fSize /= 1024
		unit++
	}
	return fmt.Sprint(
		strconv.FormatFloat(fSize, 'f', 2, 64),
		" ",
		[]string{"bytes", "KB", "MB", "GB", "TB", "PB"}[unit])
}

func calculateSize(cf ComputedFile, dir string) ComputedFile {
	f := cf.info
	if cf.size > 0 {
		return cf
	}
	if f.IsDir() {
		size, _ := dirSize(path.Join(dir, f.Name()))
		cf.size = size
	} else {
		cf.size = f.Size()
	}
	return cf
}

func main() {
	log := func(s ...interface{}) { fmt.Println(s...) }
	args := os.Args[1:]

	sortBySize := contains(args, "--sort")
	startTime := time.Now()
	dir, _ := os.Getwd()
	files, _ := ioutil.ReadDir(dir)

	var computedFiles []ComputedFile

	log(dir)

	maxLen := 0

	for _, f := range files {
		length := len(f.Name())
		if length > maxLen {
			maxLen = length + 5
		}
		computedFiles = append(computedFiles, ComputedFile{f, 0})
	}

	if sortBySize {
		sort.Slice(computedFiles, func(a, b int) bool {
			sizeA := calculateSize(computedFiles[a], dir)
			sizeB := calculateSize(computedFiles[b], dir)
			computedFiles[a] = sizeA
			computedFiles[b] = sizeB
			return sizeA.size > sizeB.size
		})
	}

	for _, f := range computedFiles {
		var calculated = calculateSize(f, dir)

		log(repeatChar(" ", 2),
			calculated.info.Name(),
			repeatChar("-", maxLen-len(calculated.info.Name())),
			prettySize(calculated.size))
	}

	defer log("\nFinished in", time.Now().Sub(startTime).String())
}
