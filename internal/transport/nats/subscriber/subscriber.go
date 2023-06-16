package subscriber

import (
	"github.com/nats-io/stan.go"
	"github.com/slavikx4/l0/internal/service"
)

type Subscriber struct {
	Connection *stan.Conn
	Service    *service.Service
}
