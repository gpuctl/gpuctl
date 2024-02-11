package config_test

import (
	"path/filepath"
	"testing"

	"github.com/gpuctl/gpuctl/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestGetSatellite_ValidConfig(t *testing.T) {
	t.Parallel()
	content := `
[groundstation]
protocol = "http"
hostname = "local.groundstation"
port = 8081

[satellite]
cache = "/tmp/satellite"
data_interval = 60
heartbeat_interval = 5`
	filename, cleanup := CreateTempConfigFile(content, t)
	defer cleanup()

	filename = filepath.Base(filename)

	config, err := config.GetSatellite(filename)
	assert.NoError(t, err)
	assert.Equal(t, "http", config.Groundstation.Protocol)
	assert.Equal(t, "local.groundstation", config.Groundstation.Hostname)
	assert.Equal(t, 8081, config.Groundstation.Port)
	assert.Equal(t, "/tmp/satellite", config.Satellite.Cache)
	assert.Equal(t, 60, config.Satellite.DataInterval)
	assert.Equal(t, 5, config.Satellite.HeartbeatInterval)
}

func TestGetSatellite_DefaultConfig(t *testing.T) {
	t.Parallel()
	content := ``

	filename, cleanup := CreateTempConfigFile(content, t)
	defer cleanup()

	filename = filepath.Base(filename)

	config, err := config.GetSatellite(filename)
	assert.NoError(t, err)
	assert.Equal(t, "http", config.Groundstation.Protocol)
	assert.Equal(t, "localhost", config.Groundstation.Hostname)
	assert.Equal(t, 8080, config.Groundstation.Port)

	assert.Equal(t, "/tmp/satellite", config.Satellite.Cache)
	assert.Equal(t, 2, config.Satellite.HeartbeatInterval)
	assert.Equal(t, 60, config.Satellite.DataInterval)
}

func TestGetSatellite_InvalidConfig(t *testing.T) {
	t.Parallel()
	content := `
groundstation: "should be a table, not a string"`

	filename, cleanup := CreateTempConfigFile(content, t)
	defer cleanup()

	filename = filepath.Base(filename)

	config, err := config.GetSatellite(filename)
	assert.Error(t, err)
	assert.Equal(t, "", config.Groundstation.Hostname)
	assert.Equal(t, 0, config.Groundstation.Port)

	assert.Equal(t, "", config.Satellite.Cache)
	assert.Equal(t, 0, config.Satellite.HeartbeatInterval)
	assert.Equal(t, 0, config.Satellite.DataInterval)
}
