package config

type ServerConfiguration struct {
	Server struct {
		Port int `toml:"port"`
	} `toml:"server"`
	Database struct {
		Url string `toml:"url"`
	} `toml:"database"`
}

func DefaultServerConfiguration() ServerConfiguration {
	return ServerConfiguration{
		Server: struct {
			Port int `toml:"port"`
		}{Port: 8080},
		Database: struct {
			Url string `toml:"url"`
		}{Url: "postgres://gpuctl@localhost/gpuctl"},
	}
}

func GetServerConfiguration(filename string) (ServerConfiguration, error) {
	return GetConfiguration[ServerConfiguration](filename, DefaultServerConfiguration)
}
