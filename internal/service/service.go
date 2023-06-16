package service

import (
	"github.com/slavikx4/l0/internal/database"
	"github.com/slavikx4/l0/internal/models"
)

type Service interface {
	AddOrder(order *models.Order) error
	GetOrder(orderUID string) (*models.Order, error)
}

type WBService struct {
	Storage *database.Storage
}

func NewWBService(storage *database.Storage) *WBService {
	return &WBService{Storage: storage}
}

//func (s *WBService) AddOrder(order *models.Order) error {
//
//	if err := s.Storage.Postgres.SetOrder(,order); err != nil{
//		return err
//	}
//}
