package config

type Server struct {
	Port int `toml:"port"`
}

type Database struct {
	Url string `toml:"url"`
}

type ControlConfiguration struct {
	Server   Server   `toml:"server"`
	Database Database `toml:"database"`
}

func (s ControlConfiguration) Merge(config Configurable) Configurable {
	sat_config, ok := config.(*ControlConfiguration)

	if !ok {
		return s
	}

	if s.Database.Url == "" {
		s.Database.Url = sat_config.Database.Url
	}

	if s.Server.Port == 0 {
		s.Server.Port = sat_config.Server.Port
	}

	return s
}

func DefaultControlConfiguration() ControlConfiguration {
	return ControlConfiguration{
		Server:   Server{Port: 8080},
		Database: Database{Url: "postgres://gpuctl@localhost/gpuctl"},
	}
}

func GetServerConfiguration(filename string) (ControlConfiguration, error) {
	return GetConfiguration[ControlConfiguration](filename, DefaultControlConfiguration)
}
