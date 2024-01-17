package main

import (
	"encoding/json"
	"fmt"
	"time"
	"github.com/gpuctl/gpuctl/internal/status/handlers"
)

func main() {
	// Just get the GPU data and print it in JSON
	gpuHandler := handlers.NvidiaGPUHandler{}
	for {
		res, err := gpuHandler.GetGPUStatus()
		if err != nil {
			fmt.Printf("Failed to parse: %v\n", err)
			return
		}

		ser, err := json.Marshal(res)
		if err != nil {
			fmt.Printf("Failed to parse: %v\n", err)
			return
		}
		fmt.Printf("%s\n", ser)
		time.Sleep(5 * time.Second)
	}
}
