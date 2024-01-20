package main

import (
	"fmt"
	"net/http"

	"github.com/gpuctl/gpuctl/cmd/groundstation/remote"
)

func main() {
	fmt.Print("Hey, world :3")

	http.HandleFunc("/api/remote", remote.HandleStatusSubmission)

}
