package main

import (
	"github.com/nats-io/stan.go"
	"github.com/slavikx4/l0/internal/transport/nats/publisher"
	er "github.com/slavikx4/l0/pkg/error"
	"github.com/slavikx4/l0/pkg/logger"
	"io"
	"os"
)

const (
	stanClusterID = "wbCluster"
	stanClientID  = "wbPublisher"
	channelName   = "orderChannel"
)

func main() {
	const op = "cmd/publisher/main -> "

	stanConnection, err := stan.Connect(stanClusterID, stanClientID)
	if err != nil {
		err = &er.Error{Err: err, Code: er.ErrorNoConnect, Message: "не удалось установить соединение с nats-streaming-server", Op: op}
		logger.Logger.Error.Fatalln(err.Error())
	}
	defer func() {
		if err := stanConnection.Close(); err != nil {
			err = &er.Error{Err: err, Code: er.ErrorClose, Message: "не удалось закрыть соединение с nats-streaming-server", Op: op}
			logger.Logger.Error.Println(err.Error())
		}
	}()

	stanPublisher := publisher.NewPublisher(&stanConnection, channelName)

	file, err := os.Open("model.json")
	if err != nil {
		err = &er.Error{Err: err, Code: er.ErrorNoConnect, Message: "не удалось открыть файл с данными model.json", Op: op}
		logger.Logger.Error.Fatalln(err.Error())
	}

	data, err := io.ReadAll(file)
	if err != nil {
		err = &er.Error{Err: err, Code: er.ErrorRead, Message: "не удалось прочитать файл", Op: op}
		logger.Logger.Error.Fatalln(err.Error())
	}

	if err := stanPublisher.Publish(channelName, data); err != nil {
		err = er.AddOp(err, op)
		logger.Logger.Error.Fatalln(err.Error())
	}
}
