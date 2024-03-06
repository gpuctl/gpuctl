package config

import (
	"time"
)

type Groundstation struct {
	Protocol string `toml:"protocol"`
	Hostname string `toml:"hostname"`
	Port     int    `toml:"port"`
}

type Satellite struct {
	Cache             string        `toml:"cache"`
	DataInterval      time.Duration `toml:"data_interval"`
	HeartbeatInterval time.Duration `toml:"heartbeat_interval"`
	FakeGPU           bool          `toml:"fake_gpu"`
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
			DataInterval:      60 * time.Second,
			HeartbeatInterval: 2 * time.Second,
			FakeGPU:           false,
		},
	}
}

func GetSatellite(filename string) (SatelliteConfiguration, error) {
	return getConfiguration[SatelliteConfiguration](filename, DefaultSatelliteConfiguration)
}
