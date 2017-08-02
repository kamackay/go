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
	http.HandleFunc("/page/", func(w http.ResponseWriter, r *http.Request) {
		defer println("Connection")
		io.WriteString(w, "<h1>This is a page!</h1>")
	})
	go println("Listening on port 7000")
	http.ListenAndServe(":7000", nil)
}
