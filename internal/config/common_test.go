package config_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gpuctl/gpuctl/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		protocol string
		hostname string
		port     int
		expected string
	}{
		{
			name:     "Standard hostname and port",
			protocol: "https",
			hostname: "example.com",
			port:     8080,
			expected: "https://example.com:8080",
		},
		{
			name:     "Localhost with common port",
			protocol: "http",
			hostname: "localhost",
			port:     8000,
			expected: "http://localhost:8000",
		},
		{
			name:     "Empty hostname over Gopher",
			protocol: "gopher",
			hostname: "",
			port:     1234,
			expected: "gopher://:1234",
		},
		{
			name:     "Zero port",
			protocol: "https",
			hostname: "example.com",
			port:     0,
			expected: "https://example.com:0",
		},
		{
			name:     "Max port number",
			protocol: "http",
			hostname: "example.com",
			port:     65535,
			expected: "http://example.com:65535",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := config.GenerateAddress(tt.protocol, tt.hostname, tt.port)
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

func TestToToml(t *testing.T) {
	t.Parallel()

	c := config.SatelliteConfiguration{
		Groundstation: config.Groundstation{"https", "foo.bar", 80},
		Satellite:     config.Satellite{"/tmp/sat", 15 * time.Second, 5 * time.Second, false},
	}

	cToml, err := config.ToToml(c)
	require.NoError(t, err)

	assert.Equal(t,
		`[groundstation]
  protocol = "https"
  hostname = "foo.bar"
  port = 80

[satellite]
  cache = "/tmp/sat"
  data_interval = "15s"
  heartbeat_interval = "5s"
  fake_gpu = false
`, cToml)

	c2 := config.ControlConfiguration{
		Timeouts: config.Timeouts{
			DeathTimeout_:    time.Second,
			MonitorInterval_: 2 * time.Second,
		},
		Server: config.Server{
			GSPort: 8080,
			WAPort: 8000,
		},
		Database: config.Database{
			InMemory:           false,
			Postgres:           true,
			PostgresUrl:        "postgres://postgres@postgres/postgres",
			DownsampleInterval: 2*time.Hour + 2*time.Minute,
		},
		Auth: config.AuthConfig{
			Username: "joe",
			Password: "mama",
		},
		SSH: config.SSHConf{
			DataDir:    "datadir",
			KeyPath:    "keypath",
			Username:   "jm",
			RemoteConf: c,
		},
	}

	c2Toml, err := config.ToToml(c2)
	require.NoError(t, err)

	assert.Equal(t,
		`[timeouts]
  death_timeout = "1s"
  monitor_interval = "2s"

[server]
  groundstation_port = 8080
  webapi_port = 8000

[database]
  inmemory = false
  postgres = true
  url = "postgres://postgres@postgres/postgres"
  downsample_interval = "2h2m0s"

[auth]
  username = "joe"
  password = "mama"

[onboard]
  datadir = "datadir"
  keyfile = "keypath"
  username = "jm"
  [onboard.remote]
    [onboard.remote.groundstation]
      protocol = "https"
      hostname = "foo.bar"
      port = 80
    [onboard.remote.satellite]
      cache = "/tmp/sat"
      data_interval = "15s"
      heartbeat_interval = "5s"
      fake_gpu = false
`, c2Toml)

}
