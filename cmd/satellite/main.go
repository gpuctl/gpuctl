package main

import (
	"fmt"
	"github.com/gpuctl/gpuctl/internal/status/handlers"
	"encoding/json"
)

func main() {
	// Just get the GPU data and print it in JSON
	gpuHandler := handlers.NvidiaGPUHandler{}
	res, err := gpuHandler.GetGPUStatus()
	if err != nil {
		fmt.Printf("Failed to parse: %v\n", err)
		return
	}

	ser, err := json.Marshal(res);
	if err != nil {
		fmt.Printf("Failed to parse: %v\n", err)
		return
	}
	fmt.Printf("%s\n", ser)
}
