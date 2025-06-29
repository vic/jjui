package jj

import (
	"github.com/BurntSushi/toml"
)

type Color struct {
	Fg        string `toml:"fg"`
	Bg        string `toml:"bg"`
	Bold      bool   `toml:"bold"`
	Underline bool   `toml:"underline"`
}

type Config struct {
	Colors        map[string]Color
	RevsetAliases map[string]string
	Revsets       struct {
		Log string
	}
}

func decodeColors(md toml.MetaData, rawColors map[string]toml.Primitive) map[string]Color {
	colorMap := make(map[string]Color)
	for name, prim := range rawColors {
		var c Color
		if err := md.PrimitiveDecode(prim, &c); err != nil {
			var fg string
			if err := md.PrimitiveDecode(prim, &fg); err == nil {
				c.Fg = fg
			}
		}
		colorMap[name] = c
	}
	return colorMap
}

type rawConfig struct {
	Colors        map[string]toml.Primitive `toml:"colors"`
	RevsetAliases map[string]string         `toml:"revset-aliases"`
	Revsets       struct {
		Log string `toml:"log"`
	} `toml:"revsets"`
}

func parseConfig(configContent string) (*Config, error) {
	var rawConfig rawConfig
	md, err := toml.Decode(configContent, &rawConfig)
	if err != nil {
		return nil, err
	}

	typedConfig := &Config{
		RevsetAliases: rawConfig.RevsetAliases,
		Revsets: struct {
			Log string
		}{
			Log: rawConfig.Revsets.Log,
		},
	}

	typedConfig.Colors = decodeColors(md, rawConfig.Colors)

	return typedConfig, nil
}

func DefaultConfig(output []byte) (*Config, error) {
	return parseConfig(string(output))
}
