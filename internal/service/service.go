package service

import (
	"github.com/slavikx4/l0/internal/database"
	"github.com/slavikx4/l0/internal/models"
	er "github.com/slavikx4/l0/pkg/error"
	"github.com/slavikx4/l0/pkg/logger"
)

// Service общий интерфейс для разных сервисов
type Service interface {
	AddOrder(order *models.Order) error
	GetOrder(orderUID string) (*models.Order, error)
}

// WBService структура для сервиса Wildberries
type WBService struct {
	Storage *database.Storage
}

func NewWBService(storage *database.Storage) *WBService {
	return &WBService{Storage: storage}
}

// AddOrder функция добавления order
func (s *WBService) AddOrder(order *models.Order) error {
	const op = "*WBService.SetOrder -> "

	if err := s.Storage.SetOrder(order); err != nil {
		err = er.AddOp(err, op)
		return err
	} else {
		logger.Logger.Process.Println("успешно записан заказ: ", order.OrderUID)
	}

	return nil
}

// GetOrder функция выдачи order
func (s *WBService) GetOrder(orderUID string) (*models.Order, error) {
	const op = "*WBService.GetOrder -> "

	order, err := s.Storage.GetOrder(orderUID)
	if err != nil {
		err = er.AddOp(err, op)
		return nil, err
	}

	return order, nil
}
