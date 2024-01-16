package main

import (
	"fmt"
	"github.com/gpuctl/gpuctl/internal/incrementer"
)

func main() {
	fmt.Println("Hello from groundstation")
	fmt.Println(incrementer.Inc(4))
}
