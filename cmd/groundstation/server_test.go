package main_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	gs "github.com/gpuctl/gpuctl/cmd/groundstation"
	"github.com/gpuctl/gpuctl/internal/uplink"
)

func TestHeartbeatRace(t *testing.T) {
	t.Parallel()

	srv := gs.NewServer()
	var wg sync.WaitGroup

	toSpawn := 100
	wg.Add(toSpawn)

	failed := false

	for i := 0; i < toSpawn; i++ {
		go func() {

			req := httptest.NewRequest(http.MethodPost, uplink.HeartbeatUrl, bytes.NewBufferString(`{}`))
			resp := httptest.NewRecorder()
			srv.ServeHTTP(resp, req)

			if resp.Code != http.StatusOK {
				failed = true
			}

			wg.Done()
		}()
	}
	wg.Wait()

	if failed {
		t.Error("one of the responces didn't return 200")
	}
}
