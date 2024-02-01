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
groundstation_port = 9090
webapi_port = 9070

[database]
postgres = true
url = "postgres://tony@ic.ac.uk/squares"
test_url = "postgres://postgres@localhost/gpuctl-tests"`
	filename, cleanup := CreateTempConfigFile(content, t)
	defer cleanup()

	filename = filepath.Base(filename)

	config, err := config.GetServerConfiguration(filename)
	assert.NoError(t, err)
	assert.Equal(t, 9090, config.Server.GSPort)
	assert.Equal(t, 9070, config.Server.WAPort)
	assert.Equal(t, true, config.Database.Postgres)
	assert.Equal(t, false, config.Database.InMemory)
	assert.Equal(t, "postgres://tony@ic.ac.uk/squares", config.Database.PostgresUrl)
	assert.Equal(t, "postgres://postgres@localhost/gpuctl-tests", config.Database.TestUrl)
}

func TestGetConfiguration_DefaultConfig(t *testing.T) {
	t.Parallel()
	content := ``

	filename, cleanup := CreateTempConfigFile(content, t)
	defer cleanup()

	filename = filepath.Base(filename)

	config, err := config.GetServerConfiguration(filename)
	assert.NoError(t, err)
	assert.Equal(t, 8080, config.Server.GSPort)
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
	assert.Equal(t, 0, config.Server.GSPort)
	assert.Equal(t, 0, config.Server.WAPort)
	assert.Equal(t, "", config.Database.PostgresUrl)
	assert.Equal(t, "", config.Database.TestUrl)
}

func TestPortToAddress(t *testing.T) {
	t.Parallel()
	assert.Equal(t, ":9090", config.PortToAddress(9090))
}

func TestControlConfigurationMerge(t *testing.T) {
	defaultConfig := config.ControlConfiguration{
		Database: config.Database{PostgresUrl: "postgres://default-url", TestUrl: "postgres://test-url"},
		Server:   config.Server{GSPort: 8080, WAPort: 8000},
	}

	fileConfig := config.ControlConfiguration{
		Database: config.Database{PostgresUrl: "", TestUrl: "postgres://non-default-test"},
		Server:   config.Server{GSPort: 9090, WAPort: 0},
	}

	mergedConfig := fileConfig.Merge(defaultConfig).(config.ControlConfiguration)

	assert.Equal(t, "postgres://default-url", mergedConfig.Database.PostgresUrl, "Expected default database URL to be applied")
	assert.Equal(t, "postgres://non-default-test", mergedConfig.Database.TestUrl, "Expected file config test URL to be applied")
	assert.Equal(t, 9090, mergedConfig.Server.GSPort, "Expected file config server groundstation port to be applied")
	assert.Equal(t, 8000, mergedConfig.Server.WAPort, "Expected default config server web api port to be applied")
}

func TestControlConfigurationMerge_UnhappyPath(t *testing.T) {
	defaultConfig := config.ControlConfiguration{}

	wrongConfig := config.SatelliteConfiguration{}

	mergedConfig := defaultConfig.Merge(&wrongConfig)

	assert.Equal(t, defaultConfig, mergedConfig, "Expected default config to be returned when wrong type is passed")
}
