package config

import (
	"bytes"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

func NewConfigFromFile(configPath string) (*Config, error) {
	viperConfig := viper.New()
	viperConfig.SetConfigType("yaml")

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	err = viperConfig.ReadConfig(bytes.NewBuffer(data))
	// Handle errors
	if err != nil {
		return nil, fmt.Errorf("fatal error config file: %w", err)
	}

	var config Config
	err = viperConfig.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	setDefaults(&config)

	return &config, nil
}

func setDefaults(config *Config) {
	if config == nil {
		return
	}

	if config.OutputDirectory == "" {
		config.OutputDirectory = "./output"
	}
	if config.DataDirectory == "" {
		config.OutputDirectory = "./data"
	}
	if len(config.DataFileExtensions) == 0 {
		config.DataFileExtensions = []string{".yaml", ".yml"}
	}
}
