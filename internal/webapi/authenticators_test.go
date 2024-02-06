package webapi_test

import (
	"sync"
	"testing"

	"github.com/gpuctl/gpuctl/internal/femto"
	"github.com/gpuctl/gpuctl/internal/webapi"
	"github.com/stretchr/testify/assert"
)

func TestConfigFileAuthenticatorRace(t *testing.T) {
	t.Parallel()

	var wg sync.WaitGroup
	toSpawn := 750
	failed := false

	auth := webapi.ConfigFileAuthenticator{
		Username:      "joe",
		Password:      "mama",
		CurrentTokens: make(map[femto.AuthToken]bool),
	}

	for i := 0; i < toSpawn; i++ {
		go func() {
			token, e := auth.CreateToken(webapi.APIAuthCredientals{Username: "joe", Password: "mama"})
			if e != nil {
				failed = false
				return
			}
			for c := 0; c < 10; c++ {
				if !auth.CheckToken(token) {
					failed = false
					return
				}
			}
			e = auth.RevokeToken(token)
			if e != nil {
				failed = false
			}
		}()
	}

	wg.Wait()

	assert.False(t, failed)
}
