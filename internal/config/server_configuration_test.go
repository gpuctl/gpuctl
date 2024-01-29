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
	filename, cleanup := config.CreateTempConfigFile(content, t)
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

	filename, cleanup := config.CreateTempConfigFile(content, t)
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
	filename, cleanup := config.CreateTempConfigFile(content, t)
	defer cleanup()

	filename = filepath.Base(filename)

	config, err := config.GetServerConfiguration(filename)
	assert.NoError(t, err)
	assert.Equal(t, 8080, config.Server.Port)
}

func TestPortToAddress(t *testing.T) {
	t.Parallel()
	assert.Equal(t, ":9090", config.PortToAddress(9090))
}
