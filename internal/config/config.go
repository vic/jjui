package config

import (
	"github.com/BurntSushi/toml"
	"os"
	"path/filepath"
)

type Config struct {
	Keys KeyMappings[keys] `toml:"keys"`
}

func findConfig() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	configPath := filepath.Join(configDir, "jjui", "config.toml")

	_, err = os.Stat(configPath)
	if err != nil {
		return "", err
	}
	return configPath, nil
}

func Load() Config {
	defaultConfig := Config{Keys: DefaultKeyMappings}
	configFile, err := findConfig()
	if err != nil {
		return defaultConfig
	}
	_, err = toml.DecodeFile(configFile, &defaultConfig)
	if err != nil {
		return defaultConfig
	}
	return defaultConfig
}
