package webapi_test

import (
	"sync"
	"testing"

	"github.com/gpuctl/gpuctl/internal/authentication"
	"github.com/gpuctl/gpuctl/internal/config"
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
		CurrentTokens: make(map[authentication.AuthToken]bool),
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

func TestAuthenticatorFromConfig(t *testing.T) {
	// A very dumb test...
	conf := config.ControlConfiguration{
		Auth: config.AuthConfig{
			Username: "user",
			Password: "password",
		},
	}
	auth := webapi.AuthenticatorFromConfig(conf)
	assert.Equal(t, conf.Auth.Username, auth.Username)
	assert.Equal(t, conf.Auth.Password, auth.Password)
}

type alwaysAuth struct{}

// CheckToken implements femto.Authenticator.
func (alwaysAuth) CheckToken(string) bool {
	return true
}

// CreateToken implements femto.Authenticator.
func (alwaysAuth) CreateToken(webapi.APIAuthCredientals) (string, error) {
	return "auth", nil
}

// RevokeToken implements femto.Authenticator.
func (alwaysAuth) RevokeToken(string) error {
	return nil
}
