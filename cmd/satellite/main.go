package main

import (
	"fmt"
	"github.com/gpuctl/gpuctl/internal/stats"
	"github.com/gpuctl/gpuctl/internal/stats/brands"
	"io/ioutil"
)

func main() {
	s := stats.GPUStatsPacket{
		Name: "gpu name",
		Brand: "some brand",
		DriverVersion: "ver 1.1",
		MemoryTotal: 1000,
		MemoryUsed: 10,
		Temp: 11.1}
	fmt.Println(s)
}
