package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

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
	return fmt.Sprint(fSize, " ", []string{"bytes", "KB", "MB", "GB", "TB", "PB"}[unit])
}

func main() {
	// args := os.Args[1:]
	startTime := time.Now()

	log := func(s ...interface{}) { fmt.Println(s...) }

	dir, _ := os.Getwd()

	log(dir)

	files, _ := ioutil.ReadDir(dir)

	maxLen := 0

	for _, f := range files {
		l := len(f.Name())
		if l > maxLen {
			maxLen = l
		}
	}

	maxLen += 5

	for _, f := range files {
		var size int64
		if f.IsDir() {
			size, _ = dirSize(path.Join(dir, f.Name()))
		} else {
			size = f.Size()
		}
		log("\t", f.Name(), repeatChar("-", maxLen-len(f.Name())), prettySize(size))
	}

	defer log("\nFinished in", time.Now().Sub(startTime).String())
}
