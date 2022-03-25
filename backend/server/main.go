package server

import (
	"net/http"

	_ "github.com/hedwig100/bookmark/backend/db"
	"github.com/hedwig100/bookmark/backend/middleware"
)

func getMux() http.ServeMux {
	var mux http.ServeMux
	mux.HandleFunc("/hello", middleware.LogWrap(hello))
	mux.HandleFunc("/users", middleware.LogWrap(postUser))
	mux.HandleFunc("/auth_test", middleware.LogWrap(middleware.Auth(hello)))
	return mux
}

func Server() http.Server {
	mux := getMux()
	server := http.Server{
		Addr:    "0.0.0.0:8081",
		Handler: &mux,
	}
	return server
}
