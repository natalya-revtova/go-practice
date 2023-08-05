package config

import (
	"errors"
)

type StorageConfig struct {
	Type string `toml:"type"`
}

func (sc StorageConfig) validate() error {
	switch sc.Type {
	case "in-memory", "sql":
		return nil
	}
	return errors.New("invalid type field")
}
