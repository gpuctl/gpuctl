package config

type Groundstation struct {
	Protocol string `toml:"protocol"`
	Hostname string `toml:"hostname"`
	Port     int    `toml:"port"`
}

type Satellite struct {
	Cache             string `toml:"cache"`
	DataInterval      int    `toml:"data_interval"`
	HeartbeatInterval int    `toml:"heartbeat_interval"`
	FakeGPU           bool   `toml:"fake_gpu"`
}

type SatelliteConfiguration struct {
	Groundstation Groundstation `toml:"groundstation"`
	Satellite     Satellite     `toml:"satellite"`
}

func DefaultSatelliteConfiguration() SatelliteConfiguration {
	return SatelliteConfiguration{
		Groundstation: Groundstation{
			Protocol: "http",
			Hostname: "localhost",
			Port:     8080,
		},
		Satellite: Satellite{
			Cache:             "/tmp/satellite",
			DataInterval:      60,
			HeartbeatInterval: 2,
			FakeGPU:           false,
		},
	}
}

func GetSatellite(filename string) (SatelliteConfiguration, error) {
	return getConfiguration[SatelliteConfiguration](filename, DefaultSatelliteConfiguration)
}