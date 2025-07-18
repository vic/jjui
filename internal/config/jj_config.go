package config

import (
	"github.com/BurntSushi/toml"
)

type JJConfig struct {
	Colors        map[string]Color  `toml:"colors"`
	RevsetAliases map[string]string `toml:"revset-aliases"`
	Revsets       struct {
		Log string `toml:"log"`
	} `toml:"revsets"`
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
