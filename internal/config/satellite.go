package config

type Groundstation struct {
	Hostname string `toml:"hostname"`
	Port     int    `toml:"port"`
}

type Satellite struct {
	Cache             string `toml:"cache"`
	DataInterval      int    `toml:"data_interval"`
	HeartbeatInterval int    `toml:"heartbeat_interval"`
}

type SatelliteConfiguration struct {
	Groundstation Groundstation `toml:"groundstation"`
	Satellite     Satellite     `toml:"satellite"`
}

func DefaultSatelliteConfiguration() SatelliteConfiguration {
	return SatelliteConfiguration{
		Groundstation: Groundstation{"localhost", 8080},
		Satellite:     Satellite{"/tmp/satellite", 60, 2},
	}
}

func GetClientConfiguration(filename string) (SatelliteConfiguration, error) {
	return getConfiguration[SatelliteConfiguration](filename, DefaultSatelliteConfiguration)
}
