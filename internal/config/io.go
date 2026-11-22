package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func getPath() (string, error) {
	path, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("Failed to get the path to $HOME directory!")
	}
	return filepath.Join(path, FILENAME), nil
}

func ReadConfig() (*Config, error) {
	path, err := getPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Failed to read the file from path \"%s\"", path)
	}

	var cfg Config

	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal byte[] to Config struct")
	}

	return &cfg, nil
}

func WriteConfig(cfg *Config) error {
	raw, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("Failed to marshal config struct")
	}
	path, err := getPath()
	if err != nil {
		return err
	}

	err = os.WriteFile(path, raw, 0o666)

	if err != nil {
		return fmt.Errorf("Failed to write file to path \"%s\"", path)
	}

	return nil
}
