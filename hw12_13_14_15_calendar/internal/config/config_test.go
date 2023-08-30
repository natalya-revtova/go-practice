package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slog"
)

func TestNewConfig(t *testing.T) {
	tests := []struct {
		description string
		path        string
		want        Config
		wantErr     bool
	}{
		{
			description: "correct config path & content",
			path:        "./testdata/valid_config.toml",
			want: Config{
				ServerHTTP: ServerHTTPConfig{
					Host:        "127.0.0.1",
					Port:        8080,
					Timeout:     4 * time.Second,
					IdleTimeout: 30 * time.Second,
				},
				ServerGRPC: ServerGRPCConfig{
					Host:              "127.0.0.1",
					Port:              50051,
					MaxConnectionIdle: 60 * time.Second,
					MaxConnectionAge:  60 * time.Second,
					Time:              60 * time.Second,
					Timeout:           60 * time.Second,
				},
				Database: DatabaseConfig{
					Host:     "127.0.0.1",
					Port:     5432,
					Username: "postgres",
					Password: "postgres",
					Database: "calendar",
				},
				Storage: StorageConfig{
					Type: "sql",
				},
				Logger: LoggerConfig{
					Level: slog.LevelInfo,
				},
			},
			wantErr: false,
		},
		{
			description: "invalid path",
			path:        "./invalid",
			want:        Config{},
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			got, err := NewConfig(tt.path)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, got, tt.want)
			}
		})
	}
}
