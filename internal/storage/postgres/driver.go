package postgres

import (
	"context"
	"embed"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lus/pasty/internal/config"
	"github.com/lus/pasty/internal/pastes"
	"github.com/lus/pasty/internal/storage"
	"github.com/rs/zerolog/log"
)

//go:embed migrations/*.sql
var migrations embed.FS

type Driver struct {
	connPool *pgxpool.Pool
	pastes   *pasteRepository
}

var _ storage.Driver = (*Driver)(nil)

func (driver *Driver) Initialize(ctx context.Context, cfg *config.Config) error {
	pool, err := pgxpool.New(ctx, cfg.Postgres.DSN)
	if err != nil {
		return err
	}

	log.Info().Msg("Performing PostgreSQL database migrations...")
	source, err := iofs.New(migrations, "migrations")
	if err != nil {
		pool.Close()
		return err
	}
	migrator, err := migrate.NewWithSourceInstance("iofs", source, cfg.Postgres.DSN)
	if err != nil {
		pool.Close()
		return err
	}
	defer func() {
		_, _ = migrator.Close()
	}()
	if err := migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		pool.Close()
		return err
	}
	log.Info().Msg("Successfully performed PostgreSQL database migrations.")

	driver.connPool = pool
	driver.pastes = &pasteRepository{
		connPool: pool,
	}
	return nil
}

func (driver *Driver) Close() error {
	driver.pastes = nil
	driver.connPool.Close()
	return nil
}

func (driver *Driver) Pastes() pastes.Repository {
	return driver.pastes
}
