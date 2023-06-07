package pastes

import (
	"context"
	"github.com/lus/pasty/internal/randx"
)

func GenerateID(ctx context.Context, repo Repository, charset string, length int) (string, error) {
	for {
		id := randx.String(charset, length)
		existing, err := repo.FindByID(ctx, id)
		if err != nil {
			return "", err
		}
		if existing == nil {
			return id, nil
		}
	}
}
