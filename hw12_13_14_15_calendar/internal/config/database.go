package config

import (
	"errors"
)

type DatabaseConfig struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	Username string `toml:"username"`
	Password string `toml:"password"`
	Database string `toml:"database"`
}

func (dc *DatabaseConfig) validate() error {
	if emptyString(dc.Host) {
		return errors.New("invalid host field")
	}
	if dc.Port <= 0 || dc.Port > 65535 {
		return errors.New("invalid port field")
	}

	if emptyString(dc.Username) {
		return errors.New("invalid username field")
	}
	if emptyString(dc.Password) {
		return errors.New("invalid password field")
	}
	if emptyString(dc.Database) {
		return errors.New("invalid database field")
	}

	return nil
}
