package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Mergable interface {
	Merge(config Mergable) Mergable
}

func PortToAddress(port int) string {
	return fmt.Sprintf(":%d", port)
}

func GenerateAddress(protocol string, hostname string, port int) string {
	return fmt.Sprintf("%s://%s%s", protocol, hostname, PortToAddress(port))
}

func IsFileEmpty(path string) (bool, error) {
	fileInfo, err := os.Stat(path)

	if err != nil {
		return false, err
	}

	return fileInfo.Size() == 0, nil
}

func getConfiguration[T Mergable](filename string, defaultGenerator func() T) (T, error) {
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

	config := defaultGenerator()

	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		// Explicit decision to zero this error state
		// as to ensure that bugs arise later if this error
		// state is misused
		var zero T
		return zero, err
	}
	return config, nil
}
