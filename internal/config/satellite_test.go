package config_test

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/BurntSushi/toml"
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
data_interval = "1m"
heartbeat_interval = "5s"`
	filename, cleanup := CreateTempConfigFile(content, t)
	defer cleanup()

	filename = filepath.Base(filename)

	conf, err := config.GetSatellite(filename)
	assert.NoError(t, err)
	assert.Equal(t, "http", conf.Groundstation.Protocol)
	assert.Equal(t, "local.groundstation", conf.Groundstation.Hostname)
	assert.Equal(t, 8081, conf.Groundstation.Port)
	assert.Equal(t, "/tmp/satellite", conf.Satellite.Cache)
	assert.Equal(t, config.Minute, conf.Satellite.DataInterval)
	assert.Equal(t, 5*config.Second, conf.Satellite.HeartbeatInterval)
}

func TestGetSatellite_DefaultConfig(t *testing.T) {
	t.Parallel()
	content := ``

	filename, cleanup := CreateTempConfigFile(content, t)
	defer cleanup()

	filename = filepath.Base(filename)

	conf, err := config.GetSatellite(filename)
	assert.NoError(t, err)
	assert.Equal(t, "http", conf.Groundstation.Protocol)
	assert.Equal(t, "localhost", conf.Groundstation.Hostname)
	assert.Equal(t, 8080, conf.Groundstation.Port)

	assert.Equal(t, "/tmp/satellite", conf.Satellite.Cache)
	assert.Equal(t, 2*config.Second, conf.Satellite.HeartbeatInterval)
	assert.Equal(t, config.Minute, conf.Satellite.DataInterval)
}

func TestGetSatellite_InvalidConfig(t *testing.T) {
	t.Parallel()
	content := `
groundstation: "should be a table, not a string"`

	filename, cleanup := CreateTempConfigFile(content, t)
	defer cleanup()

	filename = filepath.Base(filename)

	config, err := config.GetSatellite(filename)

	var parseError toml.ParseError
	if !errors.As(err, &parseError) {
		t.Fatal("Expected a parse error, but got", err)
	}

	assert.Zero(t, config)
}
