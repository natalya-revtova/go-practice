package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Logger     LoggerConfig     `toml:"logger"`
	Storage    StorageConfig    `toml:"storage"`
	ServerHTTP ServerHTTPConfig `toml:"server_http"`
	ServerGRPC ServerGRPCConfig `toml:"server_grpc"`
	Database   DatabaseConfig   `toml:"database"`
}

func NewConfig(path string) (Config, error) {
	var config Config

	_, err := toml.DecodeFile(path, &config)
	if err != nil {
		return Config{}, fmt.Errorf("parsing error: %w", err)
	}

	if err := config.validate(); err != nil {
		return Config{}, fmt.Errorf("validation error: %w", err)
	}

	return config, nil
}

func (c *Config) validate() error {
	if err := c.Logger.validate(); err != nil {
		return fmt.Errorf("invalid logger definition: %w", err)
	}
	if err := c.Storage.validate(); err != nil {
		return fmt.Errorf("invalid storage definition: %w", err)
	}
	if err := c.ServerHTTP.validate(); err != nil {
		return fmt.Errorf("invalid server_http definition: %w", err)
	}
	if err := c.ServerGRPC.validate(); err != nil {
		return fmt.Errorf("invalid server_grpc definition: %w", err)
	}
	if err := c.Database.validate(); err != nil {
		return fmt.Errorf("invalid database definition: %w", err)
	}

	return nil
}

func emptyString(str string) bool {
	return len(str) == 0
}
