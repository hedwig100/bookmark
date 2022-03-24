package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/hedwig100/bookmark/backend/middleware"
)

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

func main() {
	var mux http.ServeMux
	mux.HandleFunc("/hello", middleware.LogWrap(hello))

	server := http.Server{
		Addr:    "127.0.0.1:9080",
		Handler: &mux,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalln(fmt.Sprintf("server listening failed %v", err))
	}
}
