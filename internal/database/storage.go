package database

import (
	"context"
	in_memory "github.com/slavikx4/l0/internal/database/in-memory"
	"github.com/slavikx4/l0/internal/database/postgres"
	"github.com/slavikx4/l0/internal/models"
	er "github.com/slavikx4/l0/pkg/error"
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

func (s *Storage) SetOrder(order *models.Order) error {
	const op = "*Storage.SetOrder -> "

	if err := s.Postgres.SetOrder(context.TODO(), order); err != nil {
		err = er.AddOp(err, op)
		return err
	}

	if err := s.InMemory.SetOrder(order); err != nil {
		err = er.AddOp(err, op)
		return err
	}

	return nil
}
