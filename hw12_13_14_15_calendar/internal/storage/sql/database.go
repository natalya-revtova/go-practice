package sqlstorage

import (
	"database/sql"
	"fmt"
	"io/fs"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib" // database driver
	"github.com/jmoiron/sqlx"
	"github.com/natalya-revtova/go-practice/hw12_13_14_15_calendar/internal/app"
	"github.com/natalya-revtova/go-practice/hw12_13_14_15_calendar/internal/config"
	"github.com/pressly/goose/v3"
)

const (
	migrationsDirectory = "."
	driverName          = "pgx"
	dsnFormat           = "postgres://%s:%s@%s:%d/%s"
)

func NewConnection(config *config.DatabaseConfig) (*sqlx.DB, error) {
	db, err := sqlx.Connect(driverName, getDSN(config))
	if err != nil {
		return nil, err
	}
	return db, nil
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

type GooseLogger struct {
	app.Logger
}

func (gl *GooseLogger) Fatalf(format string, v ...interface{}) {
	gl.Error(fmt.Sprintf(format, v...))
	os.Exit(1)
}
func (gl *GooseLogger) Printf(format string, v ...interface{}) { gl.Info(fmt.Sprintf(format, v...)) }

func RunMigrations(conn *sql.DB, migrationFiles fs.FS, logger app.Logger) error {
	goose.SetBaseFS(migrationFiles)
	goose.SetLogger(&GooseLogger{logger})
	return goose.Up(conn, migrationsDirectory)
}
