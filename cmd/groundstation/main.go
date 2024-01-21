package main

import (
	"fmt"

	"github.com/gpuctl/gpuctl/internal/incrementer"
)

func main() {
	fmt.Println("Hello from ground station")
	fmt.Println(incrementer.Inc(4))
}
