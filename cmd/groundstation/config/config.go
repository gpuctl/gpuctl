package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Configuration struct {
	Server struct {
		Port int `toml:"port"`
	} `toml:"server"`
}

func DefaultConfiguration() Configuration {
	return Configuration{
		Server: struct {
			Port int `toml:"port"`
		}{Port: 8080},
	}
}

func PortToAddress(port int) string {
	return fmt.Sprintf(":%d", port)
}

func GetConfigurationFromPath(configPath string) (Configuration, error) {
	var config Configuration

	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		return DefaultConfiguration(), err
	}
	return config, nil
}

func GetConfiguration() (Configuration, error) {
	exePath, err := os.Executable()

	if err != nil {
		return DefaultConfiguration(), err
	}

	exeDir := filepath.Dir(exePath)
	configPath := filepath.Join(exeDir, "config.toml")

	return GetConfigurationFromPath(configPath)
}
