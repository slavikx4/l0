package in_memory

import (
	"context"
	"github.com/slavikx4/l0/internal/database/postgres"
	"github.com/slavikx4/l0/internal/models"
	er "github.com/slavikx4/l0/pkg/error"
	"sync"
)

type InMemory struct {
	orderUIDToOrder map[string]*models.Order
	mu              sync.RWMutex
}

func NewInMemory(ctx context.Context, postgresDB *postgres.Postgres) (*InMemory, error) {
	const op = "NewInMemory - >"

	inMemory := InMemory{
		orderUIDToOrder: make(map[string]*models.Order),
		mu:              sync.RWMutex{},
	}

	//выгружаем данные из postgres
	orders, err := postgresDB.GetOrders(ctx)
	if err != nil {
		return nil, er.AddOp(err, op)
	}

	//загружаем данные в кеш
	inMemory.SetOrders(orders)

	return &inMemory, nil
}

// SetOrders функция загружающая orders в кеш
func (m *InMemory) SetOrders(orders *[]*models.Order) {

	for _, order := range *orders {
		m.SetOrder(order)
	}
}

// SetOrder функция загружающая единственный order в кеш
func (m *InMemory) SetOrder(order *models.Order) {

	m.mu.Lock()
	m.orderUIDToOrder[order.OrderUID] = order
	m.mu.Unlock()
}

// GetOrder функция выдающая единстенный order
func (m *InMemory) GetOrder(orderUID string) (*models.Order, error) {
	const op = "*inMemory.GetOrder ->"

	m.mu.RLock()
	order, ok := m.orderUIDToOrder[orderUID]
	if !ok {
		return nil, &er.Error{Err: nil, Code: er.ErrorNotFound, Message: "order с таким orderUID не найден", Op: op}
	}

	return order, nil
}
