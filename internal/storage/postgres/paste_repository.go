package postgres

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lus/pasty/internal/pastes"
	"time"
)

type pasteRepository struct {
	connPool *pgxpool.Pool
}

var _ pastes.Repository = (*pasteRepository)(nil)

func (repo *pasteRepository) ListIDs(ctx context.Context) ([]string, error) {
	rows, _ := repo.connPool.Query(ctx, "SELECT id FROM pastes")
	result, err := pgx.CollectRows(rows, pgx.RowTo[string])
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}
	return result, nil
}

func (repo *pasteRepository) FindByID(ctx context.Context, id string) (*pastes.Paste, error) {
	rows, _ := repo.connPool.Query(ctx, "SELECT * FROM pastes WHERE id = $1", id)
	result, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByPos[pastes.Paste])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}

func (repo *pasteRepository) Upsert(ctx context.Context, paste *pastes.Paste) error {
	const query = `
		INSERT INTO pastes
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id) DO UPDATE
			SET content = excluded.content,
				"modificationToken" = excluded."modificationToken",
				metadata = excluded.metadata
	`
	_, err := repo.connPool.Exec(ctx, query, paste.ID, paste.Content, paste.ModificationToken, paste.Created, paste.Metadata)
	return err
}

func (repo *pasteRepository) DeleteByID(ctx context.Context, id string) error {
	_, err := repo.connPool.Exec(ctx, "DELETE FROM pastes WHERE id = $1", id)
	return err
}

func (repo *pasteRepository) DeleteOlderThan(ctx context.Context, age time.Duration) (int, error) {
	tag, err := repo.connPool.Exec(ctx, "DELETE FROM pastes WHERE created < $1", time.Now().Add(-age).Unix())
	if err != nil {
		return 0, err
	}
	return int(tag.RowsAffected()), nil
}
