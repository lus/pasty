package storage

import (
	"github.com/Lukaesebrot/pasty/internal/pastes"
	"github.com/bwmarrin/snowflake"
)

// Storage represents a storage type
type Storage interface {
	initialize() error
	get(id snowflake.ID) (*pastes.Paste, error)
	save(paste *pastes.Paste) error
	delete(id snowflake.ID) error
}
