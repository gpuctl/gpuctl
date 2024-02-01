package config

type Server struct {
	GSPort int `toml:"groundstation_port"`
	WAPort int `toml:"webapi_port"`
}

type Database struct {
	InMemory    bool   `toml:"inmemory"`
	Postgres    bool   `toml:"postgres"`
	PostgresUrl string `toml:"url"`
	// TODO: this is currently unused - need to figure out how to pass toml config to tests
	TestUrl string `toml:"test_url"`
}

type ControlConfiguration struct {
	Server   Server   `toml:"server"`
	Database Database `toml:"database"`
}

func (s ControlConfiguration) Merge(config Mergable) Mergable {
	sat_config, ok := config.(ControlConfiguration)

	if !ok {
		return s
	}

	if s.Database.PostgresUrl == "" {
		s.Database.PostgresUrl = sat_config.Database.PostgresUrl
	}
	if s.Database.TestUrl == "" {
		s.Database.TestUrl = sat_config.Database.TestUrl
	}

	if s.Server.GSPort == 0 {
		s.Server.GSPort = sat_config.Server.GSPort
	}
	if s.Server.WAPort == 0 {
		s.Server.WAPort = sat_config.Server.WAPort
	}

	return s
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
			TestUrl:     "postgres://postgres@localhost/gpuctl-test",
		},
	}
}

func GetServerConfiguration(filename string) (ControlConfiguration, error) {
	return getConfiguration[ControlConfiguration](filename, DefaultControlConfiguration)
}
