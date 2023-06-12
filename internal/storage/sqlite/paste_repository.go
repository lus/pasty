package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/lus/pasty/internal/pastes"
	"time"
)

type pasteRepository struct {
	connPool *sql.DB
}

var _ pastes.Repository = (*pasteRepository)(nil)

func (repo *pasteRepository) ListIDs(ctx context.Context) ([]string, error) {
	rows, err := repo.connPool.QueryContext(ctx, "SELECT id FROM pastes")
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	ids := make([]string, 0)
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, nil
}

func (repo *pasteRepository) FindByID(ctx context.Context, id string) (*pastes.Paste, error) {
	row := repo.connPool.QueryRowContext(ctx, "SELECT * FROM pastes WHERE id = ?", id)

	obj := new(pastes.Paste)

	var rawMetadata string
	if err := row.Scan(&obj.ID, &obj.Content, &obj.ModificationToken, &obj.Created, &rawMetadata); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	var metadata map[string]any
	if err := json.Unmarshal([]byte(rawMetadata), &metadata); err != nil {
		return nil, err
	}
	obj.Metadata = metadata

	return obj, nil
}

func (repo *pasteRepository) Upsert(ctx context.Context, paste *pastes.Paste) error {
	const query = `
		INSERT INTO pastes
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT (id) DO UPDATE
			SET content = excluded.content,
				modification_token = excluded.modification_token,
				metadata = excluded.metadata
	`

	rawMetadata, err := json.Marshal(paste.Metadata)
	if err != nil {
		return err
	}

	_, err = repo.connPool.ExecContext(ctx, query, paste.ID, paste.Content, paste.ModificationToken, paste.Created, rawMetadata)
	return err
}

func (repo *pasteRepository) DeleteByID(ctx context.Context, id string) error {
	_, err := repo.connPool.ExecContext(ctx, "DELETE FROM pastes WHERE id = ?", id)
	return err
}

func (repo *pasteRepository) DeleteOlderThan(ctx context.Context, age time.Duration) (int, error) {
	result, err := repo.connPool.ExecContext(ctx, "DELETE FROM pastes WHERE created < ?", time.Now().Add(-age).Unix())
	if err != nil {
		return 0, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return -1, nil
	}

	return int(affected), nil
}
