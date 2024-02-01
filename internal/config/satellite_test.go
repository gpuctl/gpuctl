package config_test

import (
	"path/filepath"
	"testing"

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
data_interval = 60
heartbeat_interval = 5`
	filename, cleanup := CreateTempConfigFile(content, t)
	defer cleanup()

	filename = filepath.Base(filename)

	config, err := config.GetClientConfiguration(filename)
	assert.NoError(t, err)
	assert.Equal(t, "local.groundstation", config.Groundstation.Hostname)
	assert.Equal(t, 8081, config.Groundstation.Port)
	assert.Equal(t, "/tmp/satellite", config.Satellite.Cache)
	assert.Equal(t, 60, config.Satellite.DataInterval)
	assert.Equal(t, 5, config.Satellite.HeartbeatInterval)
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
	assert.Equal(t, 2, config.Satellite.HeartbeatInterval)
	assert.Equal(t, 60, config.Satellite.DataInterval)
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
	assert.Equal(t, 0, config.Satellite.HeartbeatInterval)
	assert.Equal(t, 0, config.Satellite.DataInterval)
}

func TestSatelliteConfigurationMerge(t *testing.T) {
	defaultConfig := config.SatelliteConfiguration{
		Groundstation: config.Groundstation{Hostname: "localhost", Port: 8080},
		Satellite:     config.Satellite{Cache: "/tmp/default", DataInterval: 60, HeartbeatInterval: 10},
	}

	fileConfig := config.SatelliteConfiguration{
		Groundstation: config.Groundstation{Hostname: "satellite.local", Port: 0},
		Satellite:     config.Satellite{Cache: "", DataInterval: 30, HeartbeatInterval: 0},
	}

	mergedConfig := fileConfig.Merge(defaultConfig).(config.SatelliteConfiguration)

	assert.Equal(t, "satellite.local", mergedConfig.Groundstation.Hostname, "Expected file config groundstation hostname to be applied")
	assert.Equal(t, 8080, mergedConfig.Groundstation.Port, "Expected default groundstation port to be applied")
	assert.Equal(t, "/tmp/default", mergedConfig.Satellite.Cache, "Expected default satellite cache to be applied")
	assert.Equal(t, 30, mergedConfig.Satellite.DataInterval, "Expected file config satellite data interval to be applied")
	assert.Equal(t, 10, mergedConfig.Satellite.HeartbeatInterval, "Expected default satellite heartbeat interval to be applied")
}

func TestSatelliteConfigurationMerge_UnhappyPath(t *testing.T) {
	defaultConfig := config.SatelliteConfiguration{}

	wrongConfig := config.ControlConfiguration{}

	mergedConfig := defaultConfig.Merge(&wrongConfig)

	assert.Equal(t, defaultConfig, mergedConfig, "Expected default config to be returned when wrong type is passed")
}
