package config

import (
	"github.com/BurntSushi/toml"
	"os"
)

func (c *Config) Load(data string) error {
	var err error

	_, err = toml.Decode(data, c)
	if err != nil {
		return err
	}

	return nil
}

func LoadConfigFile() ([]byte, error) {
	configFile := getConfigFilePath()
	_, err := os.Stat(configFile)
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	return data, nil
}
