package database

import (
	"context"
	in_memory "github.com/slavikx4/l0/internal/database/in-memory"
	"github.com/slavikx4/l0/internal/database/postgres"
	"github.com/slavikx4/l0/internal/models"
	er "github.com/slavikx4/l0/pkg/error"
)

// Storage структура объединяющая два вида БД
type Storage struct {
	postgres *postgres.Postgres
	inMemory *in_memory.InMemory
}

func NewStorage(postgresDB *postgres.Postgres, inMemory *in_memory.InMemory) *Storage {
	return &Storage{
		postgres: postgresDB,
		inMemory: inMemory,
	}
}

// SetOrder функция устанавливающая значение в два БД
func (s *Storage) SetOrder(order *models.Order) error {
	const op = "*Storage.SetOrder -> "

	if err := s.postgres.SetOrder(context.TODO(), order); err != nil {
		err = er.AddOp(err, op)
		return err
	}

	s.inMemory.SetOrder(order)

	return nil
}

// GetOrder функция выгружающая информацию об order из кеша
func (s *Storage) GetOrder(orderUID string) (*models.Order, error) {
	const op = "*Storage.GetOrder -> "

	order, err := s.inMemory.GetOrder(orderUID)
	if err != nil {
		err = er.AddOp(err, op)
		return nil, err
	}

	return order, nil
}
