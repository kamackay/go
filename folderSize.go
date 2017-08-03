package main

import (
	"path/filepath"
	"os"
	"fmt"
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

func main() {
	args := os.Args[1:]

	log := func(s ...interface{}) { fmt.Println(s...) }

	var path string
	if len(args) > 0 {
		path = args[0]
	} else {
		log("Please input a path to get the size of")
		return
	}

	startTime := time.Now()
	log()
	log("Process size of", path)
	size, err := dirSize(path)
	if err != nil {
		panic(err)
	}
	log("\t", size, "bytes")

	var fSize float64 = float64(size)
	var unit uint8 = 0
	for fSize > 1024 {
		fSize /= 1024
		unit++
	}
	log("\t", fSize, []string{"bytes", "KB", "MB", "GB", "TB"}[unit])

	defer log("\ncalculated size of", path, "in", time.Now().Sub(startTime).String())
}
