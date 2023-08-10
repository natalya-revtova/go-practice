package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate_Database(t *testing.T) {
	config := DatabaseConfig{
		Host:     "127.0.0.1",
		Port:     8080,
		Username: "postgres",
		Password: "postgres",
		Database: "calendar",
	}

	tests := []struct {
		description string
		config      DatabaseConfig
		changeFn    func(DatabaseConfig) DatabaseConfig
		wantErr     bool
	}{
		{
			description: "valid config",
			config:      config,
			changeFn:    func(dc DatabaseConfig) DatabaseConfig { return dc },
			wantErr:     false,
		},
		{
			description: "invalid host",
			config:      config,
			changeFn: func(dc DatabaseConfig) DatabaseConfig {
				dc.Host = ""
				return dc
			},
			wantErr: true,
		},
		{
			description: "invalid port",
			config:      config,
			changeFn: func(dc DatabaseConfig) DatabaseConfig {
				dc.Port = 0
				return dc
			},
			wantErr: true,
		},
		{
			description: "invalid username",
			config:      config,
			changeFn: func(dc DatabaseConfig) DatabaseConfig {
				dc.Username = ""
				return dc
			},
			wantErr: true,
		},
		{
			description: "invalid password",
			config:      config,
			changeFn: func(dc DatabaseConfig) DatabaseConfig {
				dc.Password = ""
				return dc
			},
			wantErr: true,
		},
		{
			description: "invalid database",
			config:      config,
			changeFn: func(dc DatabaseConfig) DatabaseConfig {
				dc.Database = ""
				return dc
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			config := tt.changeFn(tt.config)
			err := config.validate()

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
