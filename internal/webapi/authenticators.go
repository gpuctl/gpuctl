package webapi

import (
	"errors"
	"sync"

	"github.com/google/uuid"
	"github.com/gpuctl/gpuctl/internal/authentication"
	"github.com/gpuctl/gpuctl/internal/config"
)

type ConfigFileAuthenticator struct {
	Username      string
	Password      string
	CurrentTokens map[authentication.AuthToken]bool
	mu            sync.Mutex
}

func AuthenticatorFromConfig(config config.ControlConfiguration) ConfigFileAuthenticator {
	return ConfigFileAuthenticator{
		Username:      config.Auth.Username,
		Password:      config.Auth.Password,
		CurrentTokens: make(map[authentication.AuthToken]bool),
	}
}

func (auth *ConfigFileAuthenticator) CreateToken(packet APIAuthCredientals) (authentication.AuthToken, error) {
	username := packet.Username
	password := packet.Password

	auth.mu.Lock()
	defer auth.mu.Unlock()

	// TODO write a proper authentication thingy
	if username != auth.Username || password != auth.Password {
		return "", authentication.InvalidCredentialsError
	}
	token := uuid.New().String()
	auth.CurrentTokens[token] = true
	return token, nil
}

func (auth *ConfigFileAuthenticator) RevokeToken(token authentication.AuthToken) error {
	auth.mu.Lock()
	defer auth.mu.Unlock()

	auth.CurrentTokens[token] = false
	return nil
}

func (auth *ConfigFileAuthenticator) CheckToken(token authentication.AuthToken) (authentication.Username, error) {
	auth.mu.Lock()
	defer auth.mu.Unlock()

	if !auth.CurrentTokens[token] {
		return "", errors.New("Bad token!")
	}

	return auth.Username, nil
}
