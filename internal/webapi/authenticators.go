package webapi

import (
	"sync"

	"github.com/google/uuid"
	"github.com/gpuctl/gpuctl/internal/config"
	"github.com/gpuctl/gpuctl/internal/femto"
)

type ConfigFileAuthenticator struct {
	Username      string
	Password      string
	CurrentTokens map[femto.AuthToken]bool
	mu            sync.Mutex
}

func AuthenticatorFromConfig(config config.ControlConfiguration) ConfigFileAuthenticator {
	return ConfigFileAuthenticator{
		Username:      config.Auth.Username,
		Password:      config.Auth.Password,
		CurrentTokens: make(map[femto.AuthToken]bool),
	}
}

func (auth *ConfigFileAuthenticator) CreateToken(packet APIAuthPacket) (femto.AuthToken, error) {
	username := packet.Username
	password := packet.Password

	auth.mu.Lock()
	defer auth.mu.Unlock()

	// TODO write a proper authentication thingy
	if username != auth.Username || password != auth.Password {
		return "", femto.InvalidCredientalsError
	}
	token := uuid.New().String()
	auth.CurrentTokens[token] = true
	return token, nil
}

func (auth *ConfigFileAuthenticator) RevokeToken(token femto.AuthToken) error {
	auth.mu.Lock()
	defer auth.mu.Unlock()

	auth.CurrentTokens[token] = false
	return nil
}

func (auth *ConfigFileAuthenticator) CheckToken(token femto.AuthToken) bool {
	auth.mu.Lock()
	defer auth.mu.Unlock()

	return auth.CurrentTokens[token]
}
