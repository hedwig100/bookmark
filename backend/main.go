package main

import (
	"fmt"
	"net/http"

	_ "github.com/hedwig100/bookmark/backend/db"
	"github.com/hedwig100/bookmark/backend/middleware"
	"github.com/hedwig100/bookmark/backend/slog"
)

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

func main() {
	var mux http.ServeMux
	mux.HandleFunc("/hello", middleware.LogWrap(hello))
	mux.HandleFunc("/auth_test", middleware.LogWrap(middleware.Auth(hello)))
	mux.HandleFunc("/auth", middleware.LogWrap(middleware.GenJWT))

	server := http.Server{
		Addr:    "0.0.0.0:8081",
		Handler: &mux,
	}

	if err := server.ListenAndServe(); err != nil {
		slog.Fatalf("server listening failed %v", err)
	}
}
