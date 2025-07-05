package config

import (
	"fmt"

	"github.com/creasty/defaults"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

func loadConfig(fileName string) (*viperConfig, error) {
	v := viper.New()

	// Set config type and name
	v.SetConfigType("yaml")
	v.SetConfigName(fileName)

	// Add multiple config paths to search
	v.AddConfigPath("./configs")
	v.AddConfigPath("../configs")
	v.AddConfigPath("../../configs")
	v.AddConfigPath("../../../configs")
	v.AddConfigPath("../../../../configs")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", fileName, err)
	}

	config := &viperConfig{}
	if err := v.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if err := defaults.Set(config); err != nil {
		return nil, fmt.Errorf("failed to set defaults from tags: %w", err)
	}

	if err := validator.New().Struct(config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return config, nil
}
