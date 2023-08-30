package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate_ServerHTTP(t *testing.T) {
	config := ServerHTTPConfig{
		Host: "127.0.0.1",
		Port: 8080,
	}

	tests := []struct {
		description string
		config      ServerHTTPConfig
		changeFn    func(ServerHTTPConfig) ServerHTTPConfig
		wantErr     bool
	}{
		{
			description: "valid config",
			config:      config,
			changeFn:    func(sc ServerHTTPConfig) ServerHTTPConfig { return sc },
			wantErr:     false,
		},
		{
			description: "invalid host",
			config:      config,
			changeFn: func(sc ServerHTTPConfig) ServerHTTPConfig {
				sc.Host = ""
				return sc
			},
			wantErr: true,
		},
		{
			description: "invalid port",
			config:      config,
			changeFn: func(sc ServerHTTPConfig) ServerHTTPConfig {
				sc.Port = 11000000
				return sc
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
