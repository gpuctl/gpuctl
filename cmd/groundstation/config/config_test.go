package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createTempConfigFile(content string, t *testing.T) (string, func()) {
	t.Helper()

	exePath, err := os.Executable()

	if err != nil {
		t.Fatal(err)
	}

	tmpfile, err := os.CreateTemp(filepath.Dir(exePath), "config.toml")

	if err != nil {
		t.Fatal(err)
	}

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}

	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	return tmpfile.Name(), func() {
		os.Remove(tmpfile.Name())
	}
}

func TestFileEmpty_EmptyCase(t *testing.T) {
	content := ``

	filename, cleanup := createTempConfigFile(content, t)
	defer cleanup()

	isEmpty, err := IsFileEmpty(filename)

	assert.NoError(t, err)
	assert.True(t, isEmpty)
}

func TestFileEmpty_NonEmptyCase(t *testing.T) {
	content := `see, it's not empty :)`

	filename, cleanup := createTempConfigFile(content, t)
	defer cleanup()

	isEmpty, err := IsFileEmpty(filename)

	assert.NoError(t, err)
	assert.False(t, isEmpty)
}

func TestFileEmpty_InvalidCase(t *testing.T) {
	_, err := IsFileEmpty("dummy_path")

	assert.Error(t, err)
}

func TestGetConfiguration_ValidConfig(t *testing.T) {
	content := `
[server]
port = 9090

[database]
url = "postgres://tony@ic.ac.uk/squares"`
	filename, cleanup := createTempConfigFile(content, t)
	defer cleanup()

	filename = filepath.Base(filename)

	config, err := GetConfiguration(filename)
	assert.NoError(t, err)
	assert.Equal(t, 9090, config.Server.Port)
	assert.Equal(t, "postgres://tony@ic.ac.uk/squares", config.Database.Url)
}

func TestGetConfiguration_DefaultConfig(t *testing.T) {
	content := ``

	filename, cleanup := createTempConfigFile(content, t)
	defer cleanup()

	filename = filepath.Base(filename)

	config, err := GetConfiguration(filename)
	assert.NoError(t, err)
	assert.Equal(t, 8080, config.Server.Port)
}

func TestGetConfiguration_InvalidConfig(t *testing.T) {
	content := `
server: "should be a table, not a string"`
	filename, cleanup := createTempConfigFile(content, t)
	defer cleanup()

	filename = filepath.Base(filename)

	config, err := GetConfiguration(filename)
	assert.NoError(t, err)
	assert.Equal(t, 8080, config.Server.Port)
}

func TestPortToAddress(t *testing.T) {
	assert.Equal(t, ":9090", PortToAddress(9090))
}
