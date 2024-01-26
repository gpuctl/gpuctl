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
	Database struct {
		Url string `toml:"url"`
	} `toml:"database"`
}

func DefaultConfiguration() Configuration {
	return Configuration{
		Server: struct {
			Port int `toml:"port"`
		}{Port: 8080},
		Database: struct {
			Url string `toml:"url"`
		}{Url: "postgres://gpuctl@localhost/gpuctl"},
	}
}

func PortToAddress(port int) string {
	return fmt.Sprintf(":%d", port)
}

func IsFileEmpty(path string) (bool, error) {
	fileInfo, err := os.Stat(path)

	if err != nil {
		return false, err
	}

	return fileInfo.Size() == 0, nil
}

func GetConfiguration(filename string) (Configuration, error) {
	exePath, err := os.Executable()

	if err != nil {
		return DefaultConfiguration(), err
	}

	exeDir := filepath.Dir(exePath)
	configPath := filepath.Join(exeDir, filename)

	fileEmpty, err := IsFileEmpty(configPath)

	if err != nil {
		return DefaultConfiguration(), err
	}

	if fileEmpty {
		return DefaultConfiguration(), nil
	}

	var config Configuration

	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		return DefaultConfiguration(), nil
	}
	return config, nil
}
