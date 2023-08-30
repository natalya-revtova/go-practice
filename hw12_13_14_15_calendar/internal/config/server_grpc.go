package config

import (
	"errors"
	"time"
)

type ServerGRPCConfig struct {
	Host              string        `toml:"host"`
	Port              int           `toml:"port"`
	MaxConnectionIdle time.Duration `toml:"max_connection_idle"`
	MaxConnectionAge  time.Duration `toml:"max_connection_age"`
	Time              time.Duration `toml:"time"`
	Timeout           time.Duration `toml:"timeout"`
}

func (sc ServerGRPCConfig) validate() error {
	if emptyString(sc.Host) {
		return errors.New("invalid host field")
	}
	if sc.Port <= 0 || sc.Port > 65535 {
		return errors.New("invalid port field")
	}
	return nil
}
