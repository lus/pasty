package storage

import (
	"context"
	"github.com/lus/pasty/internal/pastes"
)

type Driver interface {
	Initialize(ctx context.Context) error
	Close() error
	Pastes() pastes.Repository
}
