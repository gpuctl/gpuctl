package config

type ClientConfiguration struct {
	Groundstation struct {
		Hostname string `toml:"hostname"`
		Port     int    `toml:"port"`
	} `toml:"satellite"`
}

func DefaultClientConfiguration() ClientConfiguration {
	return ClientConfiguration{
		Groundstation: struct {
			Hostname string "toml:\"hostname\""
			Port     int    "toml:\"port\""
		}{"localhost", 8080},
	}
}

func GetClientConfiguration(filename string) (ClientConfiguration, error) {
	return GetConfiguration[ClientConfiguration](filename, DefaultClientConfiguration)
}
