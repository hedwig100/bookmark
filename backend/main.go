package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/hedwig100/bookmark/backend/db"
	"github.com/hedwig100/bookmark/backend/middleware"
)

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

func main() {
	var mux http.ServeMux
	mux.HandleFunc("/hello", middleware.LogWrap(hello))
	mux.HandleFunc("/auth", middleware.LogWrap(middleware.Auth(hello)))

	server := http.Server{
		Addr:    "0.0.0.0:8081",
		Handler: &mux,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalln(fmt.Sprintf("server listening failed %v", err))
	}
}
