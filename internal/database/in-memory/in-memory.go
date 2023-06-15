package in_memory

import (
	"github.com/slavikx4/l0/internal/models"
	er "github.com/slavikx4/l0/pkg/error"
	"sync"
)

type InMemory struct {
	OrderUIDToOrder map[string]*models.Order
	mu              sync.RWMutex
}

func NewInMemory() *InMemory {
	return &InMemory{
		OrderUIDToOrder: make(map[string]*models.Order),
		mu:              sync.RWMutex{},
	}
}

func (m *InMemory) SetOrder(orders *[]*models.Order) {

	for _, order := range *orders {
		m.mu.Lock()
		m.OrderUIDToOrder[order.OrderUID] = order
		m.mu.Unlock()
	}
}

func (m *InMemory) GetOrder(orderUID string) (*models.Order, error) {
	const op = "InMemory.GetOrder ->"

	m.mu.RLock()
	order, ok := m.OrderUIDToOrder[orderUID]
	if !ok {
		return nil, &er.Error{Err: nil, Code: er.ErrorNotFound, Message: "order с таким orderUID не найден", Op: op}
	}

	return order, nil
}
