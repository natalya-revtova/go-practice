package sqlstorage

import (
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib" // database driver
	"github.com/jmoiron/sqlx"
	"github.com/natalya-revtova/go-practice/hw12_13_14_15_calendar/internal/config"
)

const (
	driverName = "pgx"
	dsnFormat  = "postgres://%s:%s@%s:%d/%s"
)

func NewConnection(config *config.DatabaseConfig) (*sqlx.DB, error) {
	return sqlx.Connect(driverName, getDSN(config))
}

func getDSN(cfg *config.DatabaseConfig) string {
	return fmt.Sprintf(
		dsnFormat,
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)
}
