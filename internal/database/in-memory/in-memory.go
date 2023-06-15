package in_memory

import "github.com/slavikx4/l0/internal/models"

type InMemory struct {
	IDToOrder map[int]models.Order
	//TODO define mutex
}

func NewInMemory() *InMemory {
	return &InMemory{
		IDToOrder: make(map[int]models.Order),
	}
}
