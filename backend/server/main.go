package server

import (
	"net/http"

	"github.com/dimfeld/httptreemux/v5"
	"github.com/hedwig100/bookmark/backend/data"
	"github.com/hedwig100/bookmark/backend/middleware"
)

// NOTE: db must be initialized before the server starts to listen.
var Db data.Db

func GetMux() http.Handler {
	mux := httptreemux.NewContextMux()
	mux.GET("/hello", middleware.LogWrap(hello))
	mux.POST("/users", middleware.LogWrap(postUser))
	mux.GET("/auth_test", middleware.LogWrap(middleware.Auth(hello)))
	mux.POST("/users/:username/books", middleware.LogWrap(read))
	return mux
}

func Server() http.Server {
	Db = data.NewDbReal()
	mux := GetMux()
	server := http.Server{
		Addr:    "0.0.0.0:8081",
		Handler: mux,
	}
	return server
}
