package api

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

const GPU_STATS string = `[
{ "clock_speed": 0, "util": 0, "gpu_mem": 0, "gpu_mem_used": 0 }, 
{ "clock_speed": 0, "util": 0, "gpu_mem": 0, "gpu_mem_used": 0 },
{ "clock_speed": 0, "util": 0, "gpu_mem": 0, "gpu_mem_used": 0 }
]`

const API_URL string = "localhost:8000"

func Main(wg *sync.WaitGroup) {
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
		fmt.Fprint(w, GPU_STATS)
	default:
		{

		}
	}
}
