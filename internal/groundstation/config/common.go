package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

func PortToAddress(port int) string {
	return fmt.Sprintf(":%d", port)
}

func GenerateAddress(hostname string, port int) string {
	return fmt.Sprintf("http://%s%s", hostname, PortToAddress(port))
}

func IsFileEmpty(path string) (bool, error) {
	fileInfo, err := os.Stat(path)

	if err != nil {
		return false, err
	}

	return fileInfo.Size() == 0, nil
}

func GetConfiguration[T any](filename string, defaultGenerator func() T) (T, error) {
	exePath, err := os.Executable()

	if err != nil {
		return defaultGenerator(), err
	}

	exeDir := filepath.Dir(exePath)
	configPath := filepath.Join(exeDir, filename)

	fileEmpty, err := IsFileEmpty(configPath)

	if err != nil {
		return defaultGenerator(), err
	}

	if fileEmpty {
		return defaultGenerator(), nil
	}

	var config T

	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		return defaultGenerator(), nil
	}
	return config, nil
}
