package database

import (
	in_memory "github.com/slavikx4/l0/internal/database/in-memory"
	"github.com/slavikx4/l0/internal/database/postgres"
)

type Storage struct {
	Postgres *postgres.Postgres
	InMemory *in_memory.InMemory
}

func NewStorage(postgresDB *postgres.Postgres, inMemory *in_memory.InMemory) *Storage {
	return &Storage{
		Postgres: postgresDB,
		InMemory: inMemory,
	}
}
