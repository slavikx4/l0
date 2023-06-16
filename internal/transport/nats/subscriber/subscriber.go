package subscriber

import (
	"github.com/nats-io/stan.go"
	"github.com/slavikx4/l0/internal/service"
)

type Subscriber struct {
	Connect      stan.Conn
	DurableName  string
	Subscription stan.Subscription
	Service      *service.Service
}

func NewSubscriber(connect *stan.Conn, service service.Service, durableName string) *Subscriber {
	return &Subscriber{
		Connect:     *connect,
		DurableName: durableName,
		Service:     &service,
	}
}

func (s *Subscriber) Subscribe(channelName string) {
	//TODO дописать funс , обработать ошибку
	s.Subscription, _ = s.Connect.Subscribe(channelName, func(msg *stan.Msg) {}, stan.DurableName(s.DurableName))
}
