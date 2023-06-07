package postgres

import (
	"context"
	"embed"
	"errors"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/johejo/golang-migrate-extra/source/iofs"
	"github.com/lus/pasty/internal/config"
	"github.com/lus/pasty/internal/paste"
)

//go:embed migrations/*.sql
var migrations embed.FS

// Driver represents the Postgres storage driver
type Driver struct {
	pool               *pgxpool.Pool
	autoDeleteLifetime time.Duration
}

// Initialize initializes the Postgres storage driver
func (driver *Driver) Initialize(ctx context.Context, cfg *config.Config) error {
	pool, err := pgxpool.Connect(ctx, cfg.Postgres.DSN)
	if err != nil {
		return err
	}

	source, err := iofs.New(migrations, "migrations")
	if err != nil {
		return err
	}

	migrator, err := migrate.NewWithSourceInstance("iofs", source, cfg.Postgres.DSN)
	if err != nil {
		return err
	}

	if err := migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	driver.pool = pool
	driver.autoDeleteLifetime = cfg.AutoDelete.Lifetime
	return nil
}

// Close terminates the Postgres storage driver
func (driver *Driver) Close() error {
	driver.pool.Close()
	return nil
}

// ListIDs returns a list of all existing paste IDs
func (driver *Driver) ListIDs() ([]string, error) {
	query := "SELECT id FROM pastes"

	rows, err := driver.pool.Query(context.Background(), query)
	if err != nil {
		return []string{}, err
	}

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return []string{}, err
		}
		ids = append(ids, id)
	}

	return ids, nil
}

// Get loads a paste
func (driver *Driver) Get(id string) (*paste.Paste, error) {
	query := "SELECT * FROM pastes WHERE id = $1"

	row := driver.pool.QueryRow(context.Background(), query, id)

	paste := new(paste.Paste)
	if err := row.Scan(&paste.ID, &paste.Content, &paste.ModificationToken, &paste.Created, &paste.Metadata); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return paste, nil
}

// Save saves a paste
func (driver *Driver) Save(paste *paste.Paste) error {
	query := `
		INSERT INTO pastes (id, content, "modificationToken", created, metadata)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id) DO UPDATE
			SET content = excluded.content,
				"modificationToken" = excluded."modificationToken",
				created = excluded.created,
				metadata = excluded.metadata
	`

	_, err := driver.pool.Exec(context.Background(), query, paste.ID, paste.Content, paste.ModificationToken, paste.Created, paste.Metadata)
	return err
}

// Delete deletes a paste
func (driver *Driver) Delete(id string) error {
	query := "DELETE FROM pastes WHERE id = $1"

	_, err := driver.pool.Exec(context.Background(), query, id)
	return err
}

// Cleanup cleans up the expired pastes
func (driver *Driver) Cleanup() (int, error) {
	query := "DELETE FROM pastes WHERE created < $1"

	tag, err := driver.pool.Exec(context.Background(), query, time.Now().Add(-driver.autoDeleteLifetime).Unix())
	if err != nil {
		return 0, err
	}
	return int(tag.RowsAffected()), nil
}
