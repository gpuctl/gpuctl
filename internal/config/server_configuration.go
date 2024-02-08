package config

type Server struct {
	GSPort int `toml:"groundstation_port"`
	WAPort int `toml:"webapi_port"`
}

type Database struct {
	InMemory    bool   `toml:"inmemory"`
	Postgres    bool   `toml:"postgres"`
	PostgresUrl string `toml:"url"`
}

type AuthConfig struct {
	Username string `toml:"username"`
	Password string `toml:"password"`
}

type ControlConfiguration struct {
	Server   Server      `toml:"server"`
	Database Database    `toml:"database"`
	Auth     AuthConfig  `toml:"auth"`
	Onboard  OnboardConf `toml:"onboard"`
}

type OnboardConf struct {
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
		Onboard: OnboardConf{
			// We don't set any of the others.
			RemoteConf: DefaultSatelliteConfiguration(),
		},
	}
}

func GetServerConfiguration(filename string) (ControlConfiguration, error) {
	return getConfiguration[ControlConfiguration](filename, DefaultControlConfiguration)
}
