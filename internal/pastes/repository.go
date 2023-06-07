package pastes

import (
	"context"
	"time"
)

type Repository interface {
	ListIDs(ctx context.Context) ([]string, error)
	FindByID(ctx context.Context, id string) (*Paste, error)
	Upsert(ctx context.Context, paste *Paste) error
	DeleteByID(ctx context.Context, id string) error
	DeleteOlderThan(ctx context.Context, age time.Duration) (int, error)
}
