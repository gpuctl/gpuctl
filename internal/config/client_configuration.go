package config

type Groundstation struct {
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

func (s SatelliteConfiguration) Merge(config Mergable) Mergable {
	sat_config, ok := config.(SatelliteConfiguration)

	if !ok {
		return s
	}

	if s.Groundstation.Hostname == "" {
		s.Groundstation.Hostname = sat_config.Groundstation.Hostname
	}

	if s.Groundstation.Port == 0 {
		s.Groundstation.Port = sat_config.Groundstation.Port
	}

	if s.Satellite.Cache == "" {
		s.Satellite.Cache = sat_config.Satellite.Cache
	}

	if s.Satellite.DataInterval == 0 {
		s.Satellite.DataInterval = sat_config.Satellite.DataInterval
	}

	if s.Satellite.HeartbeatInterval == 0 {
		s.Satellite.HeartbeatInterval = sat_config.Satellite.HeartbeatInterval
	}

	return s
}

func DefaultSatelliteConfiguration() SatelliteConfiguration {
	return SatelliteConfiguration{
		Groundstation: Groundstation{"localhost", 8080},
		Satellite:     Satellite{"/tmp/satellite", 60, 2, false},
	}
}

func GetClientConfiguration(filename string) (SatelliteConfiguration, error) {
	return getConfiguration[SatelliteConfiguration](filename, DefaultSatelliteConfiguration)
}
