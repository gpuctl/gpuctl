package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gpuctl/gpuctl/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestFileEmpty_EmptyCase(t *testing.T) {
	t.Parallel()
	content := ``

	filename, cleanup := CreateTempConfigFile(content, t)
	defer cleanup()

	isEmpty, err := config.IsFileEmpty(filename)

	assert.NoError(t, err)
	assert.True(t, isEmpty)
}

func TestFileEmpty_NonEmptyCase(t *testing.T) {
	t.Parallel()
	content := `see, it's not empty :)`

	filename, cleanup := CreateTempConfigFile(content, t)
	defer cleanup()

	isEmpty, err := config.IsFileEmpty(filename)

	assert.NoError(t, err)
	assert.False(t, isEmpty)
}

func TestFileEmpty_InvalidCase(t *testing.T) {
	t.Parallel()
	_, err := config.IsFileEmpty("dummy_path")

	assert.Error(t, err)
}

func TestGenerateAddress(t *testing.T) {
	tests := []struct {
		name     string
		hostname string
		port     int
		expected string
	}{
		{
			name:     "Standard hostname and port",
			hostname: "example.com",
			port:     8080,
			expected: "http://example.com:8080",
		},
		{
			name:     "Localhost with common port",
			hostname: "localhost",
			port:     8000,
			expected: "http://localhost:8000",
		},
		{
			name:     "Empty hostname",
			hostname: "",
			port:     1234,
			expected: "http://:1234",
		},
		{
			name:     "Zero port",
			hostname: "example.com",
			port:     0,
			expected: "http://example.com:0",
		},
		{
			name:     "Max port number",
			hostname: "example.com",
			port:     65535,
			expected: "http://example.com:65535",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := config.GenerateAddress(tt.hostname, tt.port)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func CreateTempConfigFile(content string, t *testing.T) (string, func()) {
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
