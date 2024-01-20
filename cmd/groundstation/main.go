package main

import (
	"fmt"
	"net/http"

	"github.com/gpuctl/gpuctl/cmd/groundstation/config"
	"github.com/gpuctl/gpuctl/cmd/groundstation/remote"
)

func main() {
	fmt.Print("Hey, world :3")

	configuration, err := config.GetConfiguration("config.toml")

	if err != nil {
		// TODO: Using logging library for auditing, fail soft
		fmt.Println("Error detected when determining configuration")
	}

	http.HandleFunc("/api/remote", remote.HandleStatusSubmission)

	http.ListenAndServe(config.PortToAddress(configuration.Server.Port), nil)
}
