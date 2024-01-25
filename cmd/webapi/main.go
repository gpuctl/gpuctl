package main

import (
	"fmt"
	"log"
	"net/http"
)

const ALL_STATS string = `[{"name":"Shared","workStations":[{"name":"Workstation 1","gpus":[{"gpu_name":"NVIDIA GeForce GT 1030","gpu_brand":"GeForce","driver_ver":"535.146.02","memory_total":2048,"memory_util":0,"gpu_util":0,"memory_used":82,"fan_speed":35,"gpu_temp":31}]},{"name":"Workstation 2","gpus":[{"gpu_name":"NVIDIA TITAN Xp","gpu_brand":"Titan","driver_ver":"535.146.02","memory_total":12288,"memory_util":0,"gpu_util":0,"memory_used":83,"fan_speed":23,"gpu_temp":32},{"gpu_name":"NVIDIA TITAN Xp","gpu_brand":"Titan","driver_ver":"535.146.02","memory_total":12288,"memory_util":0,"gpu_util":0,"memory_used":83,"fan_speed":23,"gpu_temp":32}]},{"name":"Workstation 3","gpus":[{"gpu_name":"NVIDIA GeForce GT 730","gpu_brand":"GeForce","driver_ver":"470.223.02","memory_total":2001,"memory_util":0,"gpu_util":0,"memory_used":220,"fan_speed":30,"gpu_temp":27}]},{"name":"Workstation 5","gpus":[{"gpu_name":"NVIDIA TITAN Xp","gpu_brand":"Titan","driver_ver":"535.146.02","memory_total":12288,"memory_util":0,"gpu_util":0,"memory_used":83,"fan_speed":23,"gpu_temp":32},{"gpu_name":"NVIDIA TITAN Xp","gpu_brand":"Titan","driver_ver":"535.146.02","memory_total":12288,"memory_util":0,"gpu_util":0,"memory_used":83,"fan_speed":23,"gpu_temp":32}]},{"name":"Workstation 4","gpus":[{"gpu_name":"NVIDIA GeForce GT 1030","gpu_brand":"GeForce","driver_ver":"535.146.02","memory_total":2048,"memory_util":0,"gpu_util":0,"memory_used":82,"fan_speed":35,"gpu_temp":31}]},{"name":"Workstation 6","gpus":[{"gpu_name":"NVIDIA GeForce GT 730","gpu_brand":"GeForce","driver_ver":"470.223.02","memory_total":2001,"memory_util":0,"gpu_util":0,"memory_used":220,"fan_speed":30,"gpu_temp":27}]}]}]`

const API_URL string = "localhost:8000"

func main() {
	init_server()
	for {
		// Loop
	}
}

func init_server() {
	fmt.Println("API is starting up!")
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(API_URL, nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	switch r.URL.Path {
	case "/api/stats/all":
		fmt.Fprint(w, ALL_STATS)
	default:
		{

		}
	}
}
