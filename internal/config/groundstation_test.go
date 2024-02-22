package config_test

import (
	"path/filepath"
	"testing"

	"github.com/gpuctl/gpuctl/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestGetControl_ValidConfig(t *testing.T) {
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

	config, err := config.GetControl(filename)
	assert.NoError(t, err)
	assert.Equal(t, 9090, config.Server.GSPort)
	assert.Equal(t, 9070, config.Server.WAPort)
	assert.Equal(t, true, config.Database.Postgres)
	assert.Equal(t, false, config.Database.InMemory)
	assert.Equal(t, "postgres://tony@ic.ac.uk/squares", config.Database.PostgresUrl)
}

func TestGetControl_DefaultConfig(t *testing.T) {
	t.Parallel()
	content := ``

	filename, cleanup := CreateTempConfigFile(content, t)
	defer cleanup()

	filename = filepath.Base(filename)

	config, err := config.GetControl(filename)
	assert.NoError(t, err)
	assert.Equal(t, 8080, config.Server.GSPort)
}

func TestGetControl_InvalidConfig(t *testing.T) {
	t.Parallel()
	content := `
server: "should be a table, not a string"`
	filename, cleanup := CreateTempConfigFile(content, t)
	defer cleanup()

	filename = filepath.Base(filename)

	conf, err := config.GetControl(filename)
	assert.Error(t, err)
	assert.Equal(t, config.ControlConfiguration{}, conf)
}

func TestPortToAddress(t *testing.T) {
	t.Parallel()
	assert.Equal(t, ":9090", config.PortToAddress(9090))
}
