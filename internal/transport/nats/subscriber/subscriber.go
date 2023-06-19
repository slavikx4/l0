package subscriber

import (
	"encoding/json"
	"github.com/nats-io/stan.go"
	"github.com/slavikx4/l0/internal/models"
	"github.com/slavikx4/l0/internal/service"
	er "github.com/slavikx4/l0/pkg/error"
	"github.com/slavikx4/l0/pkg/logger"
)

type Subscriber struct {
	Connect      stan.Conn
	DurableName  string
	Subscription stan.Subscription
	Service      service.Service
}

func NewSubscriber(connect *stan.Conn, service service.Service, durableName string) *Subscriber {
	return &Subscriber{
		Connect:     *connect,
		DurableName: durableName,
		Service:     service,
	}
}

func (s *Subscriber) Subscribe(channelName string) error {
	const op = "*Subscriber.Subscribe -> "

	var err error
	s.Subscription, err = s.Connect.Subscribe(channelName, s.handlMessage, stan.DurableName(s.DurableName))
	if err != nil {
		return &er.Error{Err: err, Code: er.ErrorSubscribe, Message: "не удалось подписаться на канал", Op: op}
	}

	return nil
}

func (s *Subscriber) handlMessage(msg *stan.Msg) {
	const op = "*Subscriber.handleMessage -> "

	logger.Logger.Process.Println("пришло новое сообщение через канал")

	order := models.Order{}

	if err := json.Unmarshal(msg.Data, &order); err != nil {
		err = &er.Error{Err: err, Code: er.ErrorJson, Message: "не удалось раскодировать json в order", Op: op}
		logger.Logger.Error.Println(err.Error())
	}

	logger.Logger.Process.Println("успешно раскодировано сообщение orderUID: ", order.OrderUID)

	if err := s.Service.AddOrder(&order); err != nil {
		err = er.AddOp(err, op)
		logger.Logger.Error.Println(err.Error())
	}
}
