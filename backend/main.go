package main

import (
	"fmt"
	"net/http"
)

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

func main() {
	var mux http.ServeMux
	mux.HandleFunc("/hello", hello)

	server := http.Server{
		Addr:    "127.0.0.1:9080",
		Handler: &mux,
	}
	server.ListenAndServe()
}
