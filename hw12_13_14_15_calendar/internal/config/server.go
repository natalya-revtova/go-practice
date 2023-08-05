package config

import (
	"errors"
	"time"
)

type ServerConfig struct {
	Host        string        `toml:"host"`
	Port        int           `toml:"port"`
	Timeout     time.Duration `toml:"timeout"`
	IdleTimeout time.Duration `toml:"idle_timeout"`
}

func (sc ServerConfig) validate() error {
	if emptyString(sc.Host) {
		return errors.New("invalid host field")
	}
	if sc.Port <= 0 || sc.Port > 65535 {
		return errors.New("invalid port field")
	}

	return nil
}
