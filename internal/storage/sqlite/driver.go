package sqlite

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/lus/pasty/internal/pastes"
	"github.com/lus/pasty/internal/storage"
	"github.com/rs/zerolog/log"
	_ "modernc.org/sqlite"
)

//go:embed migrations/*.sql
var migrations embed.FS

type Driver struct {
	filePath string
	connPool *sql.DB
	pastes   *pasteRepository
}

var _ storage.Driver = (*Driver)(nil)

func New(filePath string) *Driver {
	return &Driver{
		filePath: filePath,
	}
}

func (driver *Driver) Initialize(ctx context.Context) error {
	db, err := sql.Open("sqlite", driver.filePath)
	if err != nil {
		return err
	}
	if err := db.PingContext(ctx); err != nil {
		return err
	}

	log.Info().Msg("Performing SQLite database migrations...")
	source, err := iofs.New(migrations, "migrations")
	if err != nil {
		_ = db.Close()
		return err
	}
	defer func() {
		_ = source.Close()
	}()
	migrateDriver, err := sqlite.WithInstance(db, &sqlite.Config{
		MigrationsTable: sqlite.DefaultMigrationsTable,
		DatabaseName:    driver.filePath,
		NoTxWrap:        false,
	})
	if err != nil {
		_ = db.Close()
		return err
	}
	migrator, err := migrate.NewWithInstance("iofs", source, "sqlite", migrateDriver)
	if err != nil {
		_ = db.Close()
		return err
	}
	if err := migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		_ = db.Close()
		return err
	}

	driver.connPool = db
	driver.pastes = &pasteRepository{
		connPool: db,
	}
	return nil
}

func (driver *Driver) Close() error {
	driver.pastes = nil
	_ = driver.connPool.Close()
	driver.connPool = nil
	return nil
}

func (driver *Driver) Pastes() pastes.Repository {
	return driver.pastes
}
