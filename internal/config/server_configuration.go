package config

type Server struct {
	GSPort int `toml:"groundstation_port"`
	WAPort int `toml:"webapi_port"`
}

type Database struct {
	InMemory           bool   `toml:"inmemory"`
	Postgres           bool   `toml:"postgres"`
	PostgresUrl        string `toml:"url"`
	DownsampleInterval int    `toml:"downsample_interval`
}

type AuthConfig struct {
	Username string `toml:"username"`
	Password string `toml:"password"`
}

type ControlConfiguration struct {
	Server   Server     `toml:"server"`
	Database Database   `toml:"database"`
	Auth     AuthConfig `toml:"auth"`
}

func DefaultControlConfiguration() ControlConfiguration {
	return ControlConfiguration{
		Server: Server{
			GSPort: 8080,
			WAPort: 8000,
		},
		Database: Database{
			InMemory:    false,
			Postgres:    false,
			PostgresUrl: "postgres://postgres@postgres/postgres",
		},
		Auth: AuthConfig{
			Username: "admin",
			Password: "password",
		},
	}
}

func GetServerConfiguration(filename string) (ControlConfiguration, error) {
	return getConfiguration[ControlConfiguration](filename, DefaultControlConfiguration)
}
