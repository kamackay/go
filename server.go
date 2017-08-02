package main

import (
	"io"
	"net/http"
	"fmt"
)

func main() {
	const port = 7000

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer println("Connection")
		io.WriteString(w, "<h1>Hello world</h1><h2>Ver2</h2>")
	})
	http.HandleFunc("/page/", func(w http.ResponseWriter, r *http.Request) {
		defer println("Connection")
		io.WriteString(w, "<h1>This is a page!</h1>")
	})
	go println(fmt.Sprintf("Listening on port %d", port))
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
