package main

import (
	"io"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer println("Connection")
		io.WriteString(w, "<h1>Hello world</h1><h2>Ver2</h2>")
	})
	http.ListenAndServe(":7000", nil)
}
