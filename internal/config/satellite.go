package config

import "time"

type Groundstation struct {
	Hostname string `toml:"hostname"`
	Port     int    `toml:"port"`
}

type Satellite struct {
	Cache              string   `toml:"cache"`
	DataInterval_      Duration `toml:"data_interval"`
	HeartbeatInterval_ Duration `toml:"heartbeat_interval"`
}

func (s Satellite) DataInterval() time.Duration {
	return time.Duration(s.DataInterval_)
}
func (s Satellite) HeartbeatInterval() time.Duration {
	return time.Duration(s.HeartbeatInterval_)
}

type SatelliteConfiguration struct {
	Groundstation Groundstation `toml:"groundstation"`
	Satellite     Satellite     `toml:"satellite"`
}

func DefaultSatelliteConfiguration() SatelliteConfiguration {
	return SatelliteConfiguration{
		Groundstation: Groundstation{"localhost", 8080},
		Satellite:     Satellite{"/tmp/satellite", 60 * Second, 2 * Second},
	}
}

func GetClientConfiguration(filename string) (SatelliteConfiguration, error) {
	return getConfiguration[SatelliteConfiguration](filename, DefaultSatelliteConfiguration)
}
