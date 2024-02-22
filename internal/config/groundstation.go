package config

type Server struct {
	GSPort int `toml:"groundstation_port"`
	WAPort int `toml:"webapi_port"`
}

type Timeouts struct {
	DeathTimeout    int `toml:"death_timeout"`
	MonitorInterval int `toml:"monitor_interval"`
}

type Database struct {
	InMemory           bool   `toml:"inmemory"`
	Postgres           bool   `toml:"postgres"`
	PostgresUrl        string `toml:"url"`
	DownsampleInterval int    `toml:"downsample_interval"`
}

type AuthConfig struct {
	Username string `toml:"username"`
	Password string `toml:"password"`
}

type ControlConfiguration struct {
	Timeouts Timeouts   `toml:"timeouts"`
	Server   Server     `toml:"server"`
	Database Database   `toml:"database"`
	Auth     AuthConfig `toml:"auth"`
	SSH      SSHConf    `toml:"onboard"` // TODO: Change name to ssh_configuration, deferred due to it being a breaking change
}

type SSHConf struct {
	DataDir    string                 `toml:"datadir"`
	KeyPath    string                 `toml:"keyfile"`
	Username   string                 `toml:"username"`
	RemoteConf SatelliteConfiguration `toml:"remote"`
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
		SSH: SSHConf{
			// We don't set any of the others.
			RemoteConf: DefaultSatelliteConfiguration(),
		},
	}
}

func GetControl(filename string) (ControlConfiguration, error) {
	return getConfiguration[ControlConfiguration](filename, DefaultControlConfiguration)
}
