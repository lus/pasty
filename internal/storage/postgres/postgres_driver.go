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
	"github.com/lus/pasty/internal/shared"
)

//go:embed migrations/*.sql
var migrations embed.FS

// PostgresDriver represents the Postgres storage driver
type PostgresDriver struct {
	pool *pgxpool.Pool
}

// Initialize initializes the Postgres storage driver
func (driver *PostgresDriver) Initialize() error {
	pool, err := pgxpool.Connect(context.Background(), config.Current.Postgres.DSN)
	if err != nil {
		return err
	}

	source, err := iofs.New(migrations, "migrations")
	if err != nil {
		return err
	}

	migrator, err := migrate.NewWithSourceInstance("iofs", source, config.Current.Postgres.DSN)
	if err != nil {
		return err
	}

	if err := migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	driver.pool = pool
	return nil
}

// Terminate terminates the Postgres storage driver
func (driver *PostgresDriver) Terminate() error {
	driver.pool.Close()
	return nil
}

// ListIDs returns a list of all existing paste IDs
func (driver *PostgresDriver) ListIDs() ([]string, error) {
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
func (driver *PostgresDriver) Get(id string) (*shared.Paste, error) {
	query := "SELECT * FROM pastes WHERE id = $1"

	row := driver.pool.QueryRow(context.Background(), query, id)

	paste := new(shared.Paste)
	if err := row.Scan(&paste.ID, &paste.Content, &paste.ModificationToken, &paste.Created, &paste.AutoDelete, &paste.Metadata); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return paste, nil
}

// Save saves a paste
func (driver *PostgresDriver) Save(paste *shared.Paste) error {
	query := `
		INSERT INTO pastes (id, content, modificationToken, created, autoDelete)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id) DO UPDATE
			SET content = excluded.token,
				modificationToken = excluded.modificationToken,
				created = excluded.created,
				autoDelete = excluded.autoDelete,
				metadata = excluded.metadata
	`

	_, err := driver.pool.Exec(context.Background(), query, paste.ID, paste.Content, paste.ModificationToken, paste.Created, paste.AutoDelete, paste.Metadata)
	return err
}

// Delete deletes a paste
func (driver *PostgresDriver) Delete(id string) error {
	query := "DELETE FROM pastes WHERE id = $1"

	_, err := driver.pool.Exec(context.Background(), query, id)
	return err
}

// Cleanup cleans up the expired pastes
func (driver *PostgresDriver) Cleanup() (int, error) {
	query := "DELETE FROM pastes WHERE autoDelete = true AND created < $2"

	tag, err := driver.pool.Exec(context.Background(), query, time.Now().Add(-config.Current.AutoDelete.Lifetime).Unix())
	if err != nil {
		return 0, err
	}
	return int(tag.RowsAffected()), nil
}
