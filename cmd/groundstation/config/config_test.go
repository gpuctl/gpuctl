package config

import (
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/stretchr/testify/assert"
)

func TestGetConfigurationFromPath(t *testing.T) {
	const validConfig = `
[server]
port = 9090`

	var config Configuration
	_, err := toml.Decode(validConfig, &config)

	assert.NoError(t, err)
	assert.Equal(t, 9090, config.Server.Port)
}

func TestPortToAddress(t *testing.T) {
	assert.Equal(t, ":9090", PortToAddress(9090))
}
