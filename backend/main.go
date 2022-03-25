package main

import (
	"github.com/hedwig100/bookmark/backend/server"
	"github.com/hedwig100/bookmark/backend/slog"
)

func main() {
	server := server.Server()
	if err := server.ListenAndServe(); err != nil {
		slog.Fatalf("server listening failed %v", err)
	}
}
