package config_test

import (
	"path/filepath"
	"testing"

	"github.com/gpuctl/gpuctl/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestGetConfiguration_ValidConfig(t *testing.T) {
	t.Parallel()
	content := `
[server]
port = 9090

[database]
url = "postgres://tony@ic.ac.uk/squares"`
	filename, cleanup := CreateTempConfigFile(content, t)
	defer cleanup()

	filename = filepath.Base(filename)

	config, err := config.GetServerConfiguration(filename)
	assert.NoError(t, err)
	assert.Equal(t, 9090, config.Server.Port)
	assert.Equal(t, "postgres://tony@ic.ac.uk/squares", config.Database.Url)
}

func TestGetConfiguration_DefaultConfig(t *testing.T) {
	t.Parallel()
	content := ``

	filename, cleanup := CreateTempConfigFile(content, t)
	defer cleanup()

	filename = filepath.Base(filename)

	config, err := config.GetServerConfiguration(filename)
	assert.NoError(t, err)
	assert.Equal(t, 8080, config.Server.Port)
}

func TestGetConfiguration_InvalidConfig(t *testing.T) {
	t.Parallel()
	content := `
server: "should be a table, not a string"`
	filename, cleanup := CreateTempConfigFile(content, t)
	defer cleanup()

	filename = filepath.Base(filename)

	config, err := config.GetServerConfiguration(filename)
	assert.Error(t, err)
	assert.Equal(t, 8080, config.Server.Port)
}

func TestPortToAddress(t *testing.T) {
	t.Parallel()
	assert.Equal(t, ":9090", config.PortToAddress(9090))
}

func TestControlConfigurationMerge(t *testing.T) {
	defaultConfig := config.ControlConfiguration{
		Database: config.Database{Url: "postgres://default-url"},
		Server:   config.Server{Port: 8080},
	}

	fileConfig := config.ControlConfiguration{
		Database: config.Database{Url: ""},
		Server:   config.Server{Port: 9090},
	}

	mergedConfig := fileConfig.Merge(defaultConfig).(config.ControlConfiguration)

	assert.Equal(t, "postgres://default-url", mergedConfig.Database.Url, "Expected default database URL to be applied")
	assert.Equal(t, 9090, mergedConfig.Server.Port, "Expected file config server port to be applied")
}

func TestControlConfigurationMerge_UnhappyPath(t *testing.T) {
	defaultConfig := config.ControlConfiguration{}

	wrongConfig := config.SatelliteConfiguration{}

	mergedConfig := defaultConfig.Merge(&wrongConfig)

	assert.Equal(t, defaultConfig, mergedConfig, "Expected default config to be returned when wrong type is passed")
}
