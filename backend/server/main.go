package server

import (
	"net/http"

	"github.com/hedwig100/bookmark/backend/data"
	"github.com/hedwig100/bookmark/backend/middleware"
)

// NOTE: db must be initialized before the server starts to listen.
var Db data.Db

func GetMux() http.ServeMux {
	var mux http.ServeMux
	mux.HandleFunc("/hello", middleware.LogWrap(hello))
	mux.HandleFunc("/users", middleware.LogWrap(postUser))
	mux.HandleFunc("/auth_test", middleware.LogWrap(middleware.Auth(hello)))
	return mux
}

func Server() http.Server {
	Db = data.NewDbReal()
	mux := GetMux()
	server := http.Server{
		Addr:    "0.0.0.0:8081",
		Handler: &mux,
	}
	return server
}
