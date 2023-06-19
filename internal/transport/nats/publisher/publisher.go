package publisher

import (
	"github.com/nats-io/stan.go"
	er "github.com/slavikx4/l0/pkg/error"
)

// Publisher структура паблишера в канал nats-streaming
type Publisher struct {
	Connection  stan.Conn
	ChannelName string
}

func NewPublisher(connection *stan.Conn, channelName string) *Publisher {
	return &Publisher{
		Connection:  *connection,
		ChannelName: channelName,
	}
}

// Publish функция для посылки сообщения в канал
func (p *Publisher) Publish(channelName string, data []byte) error {
	const op = "Publisher.Publish -> "

	if err := p.Connection.Publish(channelName, data); err != nil {
		return &er.Error{Err: err, Code: er.ErrorPublish, Message: "не удалось добавить данные в канал", Op: op}
	}
	return nil
}
