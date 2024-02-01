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

func DefaultControlConfiguration() ControlConfiguration {
	return ControlConfiguration{
		Server:   Server{Port: 8080},
		Database: Database{Url: "postgres://gpuctl@localhost/gpuctl"},
	}
}

func GetServerConfiguration(filename string) (ControlConfiguration, error) {
	return getConfiguration[ControlConfiguration](filename, DefaultControlConfiguration)
}
