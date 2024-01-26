package main

import (
	"log/slog"
	"net/http"

	"github.com/gpuctl/gpuctl/internal/webapi"
)

func main() {
	slog.Info("Starting webapi")

	// TODO: Don't be nil
	srv := webapi.NewServer(nil)

	http.ListenAndServe(":8000", srv)
}
