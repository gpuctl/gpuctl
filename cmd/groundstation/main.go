package main

import (
	"fmt"
	"sync"

	"github.com/gpuctl/gpuctl/cmd/groundstation/api"
	"github.com/gpuctl/gpuctl/internal/incrementer"
)

func main() {
	fmt.Println("Hello from ground station")
	fmt.Println(incrementer.Inc(4))
	var wg sync.WaitGroup
	wg.Add(1)
	go api.Main(&wg)
	wg.Wait()
}
