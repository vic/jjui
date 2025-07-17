package config

import (
	"github.com/BurntSushi/toml"
)

type JJConfig struct {
	Colors        map[string]Color
	RevsetAliases map[string]string
	Revsets       struct {
		Log string
	}
}

func parseConfig(configContent string) (*JJConfig, error) {
	var config JJConfig
	_, err := toml.Decode(configContent, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func DefaultConfig(output []byte) (*JJConfig, error) {
	return parseConfig(string(output))
}
