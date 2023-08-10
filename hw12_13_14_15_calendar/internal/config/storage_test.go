package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate_Storage(t *testing.T) {
	config := StorageConfig{
		Type: "sql",
	}

	tests := []struct {
		description string
		config      StorageConfig
		changeFn    func(StorageConfig) StorageConfig
		wantErr     bool
	}{
		{
			description: "valid config",
			config:      config,
			changeFn:    func(sc StorageConfig) StorageConfig { return sc },
			wantErr:     false,
		},
		{
			description: "invalid type",
			config:      config,
			changeFn: func(sc StorageConfig) StorageConfig {
				sc.Type = "invalid"
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
