package main

import (
	"log/slog"
	"net/http"

	"github.com/gpuctl/gpuctl/internal/webapi"
)

func main() {
	slog.Info("Starting webapi")

	srv := webapi.NewServer()

	http.ListenAndServe(":8000", srv)
}
