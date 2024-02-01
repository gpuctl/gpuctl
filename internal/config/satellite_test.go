package config_test

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/gpuctl/gpuctl/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestGetClientConfiguration_ValidConfig(t *testing.T) {
	t.Parallel()
	content := `
[groundstation]
hostname = "local.groundstation"
port = 8081

[satellite]
cache = "/tmp/satellite"
data_interval = "666s"
heartbeat_interval = "55s"`
	filename, cleanup := CreateTempConfigFile(content, t)
	defer cleanup()

	filename = filepath.Base(filename)

	config, err := config.GetClientConfiguration(filename)
	assert.NoError(t, err)
	assert.Equal(t, "local.groundstation", config.Groundstation.Hostname)
	assert.Equal(t, 8081, config.Groundstation.Port)
	assert.Equal(t, "/tmp/satellite", config.Satellite.Cache)
	assert.Equal(t, 666*time.Second, config.Satellite.DataInterval())
	assert.Equal(t, 55*time.Second, config.Satellite.HeartbeatInterval())
}

func TestGetClientConfiguration_DefaultConfig(t *testing.T) {
	t.Parallel()
	content := ``

	filename, cleanup := CreateTempConfigFile(content, t)
	defer cleanup()

	filename = filepath.Base(filename)

	config, err := config.GetClientConfiguration(filename)
	assert.NoError(t, err)
	assert.Equal(t, "localhost", config.Groundstation.Hostname)
	assert.Equal(t, 8080, config.Groundstation.Port)

	assert.Equal(t, "/tmp/satellite", config.Satellite.Cache)
	assert.Equal(t, 60*time.Second, config.Satellite.DataInterval())
	assert.Equal(t, 2*time.Second, config.Satellite.HeartbeatInterval())
}

func TestGetClientConfiguration_InvalidConfig(t *testing.T) {
	t.Parallel()
	content := `
groundstation: "should be a table, not a string"`

	filename, cleanup := CreateTempConfigFile(content, t)
	defer cleanup()

	filename = filepath.Base(filename)

	config, err := config.GetClientConfiguration(filename)
	assert.Error(t, err)
	assert.Equal(t, "", config.Groundstation.Hostname)
	assert.Equal(t, 0, config.Groundstation.Port)

	assert.Equal(t, "", config.Satellite.Cache)
	assert.Equal(t, time.Duration(0), config.Satellite.HeartbeatInterval())
	assert.Equal(t, time.Duration(0), config.Satellite.DataInterval())
}
