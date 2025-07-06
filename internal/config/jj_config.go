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

func parseConfig(configContent string) (*JJConfig, error) {
	type rawConfig struct {
		Colors        map[string]toml.Primitive `toml:"colors"`
		RevsetAliases map[string]string         `toml:"revset-aliases"`
		Revsets       struct {
			Log string `toml:"log"`
		} `toml:"revsets"`
	}

	var raw rawConfig
	md, err := toml.Decode(configContent, &raw)
	if err != nil {
		return nil, err
	}

	typedConfig := &JJConfig{
		RevsetAliases: raw.RevsetAliases,
		Revsets: struct {
			Log string
		}{
			Log: raw.Revsets.Log,
		},
	}

	typedConfig.Colors = decodeColors(md, raw.Colors)

	return typedConfig, nil
}

func DefaultConfig(output []byte) (*JJConfig, error) {
	return parseConfig(string(output))
}
