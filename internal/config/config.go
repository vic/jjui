package config

import (
	"github.com/BurntSushi/toml"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

type Config struct {
	Keys KeyMappings[keys] `toml:"keys"`
}

func getConfigFilePath() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return ""
	}
	return filepath.Join(configDir, "jjui", "config.toml")
}

func getDefaultEditor() string {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = os.Getenv("VISUAL")
	}

	// Fallback to common editors if not set
	if editor == "" {
		candidates := []string{"nano", "vim", "vi", "notepad.exe"} // Windows fallback
		for _, candidate := range candidates {
			if p, err := exec.LookPath(candidate); err == nil {
				editor = p
				break
			}
		}
	}

	return editor
}

func Load() Config {
	defaultConfig := Config{Keys: DefaultKeyMappings}
	configFile := getConfigFilePath()
	_, err := os.Stat(configFile)
	if err != nil {
		return defaultConfig
	}
	_, err = toml.DecodeFile(configFile, &defaultConfig)
	if err != nil {
		Current = &defaultConfig
		return defaultConfig
	}
	return defaultConfig
}

func Edit() int {
	configFile := getConfigFilePath()
	_, err := os.Stat(configFile)
	if os.IsNotExist(err) {
		configPath := path.Dir(configFile)
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			err = os.MkdirAll(configPath, 0755)
			if err != nil {
				log.Fatal(err)
				return -1
			}
		}
		if _, err := os.Stat(configFile); os.IsNotExist(err) {
			_, err := os.Create(configFile)
			if err != nil {
				log.Fatal(err)
				return -1
			}
		}
	}

	editor := getDefaultEditor()
	if editor == "" {
		log.Fatal("No editor found. Please set $EDITOR or $VISUAL")
	}

	cmd := exec.Command(editor, configFile)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	return cmd.ProcessState.ExitCode()
}
