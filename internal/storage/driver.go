package storage

import (
	"context"
	"github.com/lus/pasty/internal/config"
	"github.com/lus/pasty/internal/pastes"
)

type Driver interface {
	Initialize(ctx context.Context, cfg *config.Config) error
	Close() error
	Pastes() pastes.Repository
}
