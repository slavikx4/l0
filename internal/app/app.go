package app

import (
	"context"
	"fmt"
	"github.com/nats-io/stan.go"
	"github.com/slavikx4/l0/internal/database"
	in_memory "github.com/slavikx4/l0/internal/database/in-memory"
	"github.com/slavikx4/l0/internal/database/postgres"
	"github.com/slavikx4/l0/internal/service"
	"github.com/slavikx4/l0/internal/transport/nats/subscriber"
	"github.com/slavikx4/l0/internal/transport/rest"
	er "github.com/slavikx4/l0/pkg/error"
	"github.com/slavikx4/l0/pkg/logger"
	"github.com/spf13/viper"
	"net/http"
)

const (
	portKey                  = "port"
	dbKey                    = "db"
	stanClusterIDKey         = "stanClusterID"
	stanClientIDKey          = "stanClientID"
	subscriberDurableNameKey = "subscriberDurableName"
	channelNameKey           = "channelName"
)

func Run() {
	const op = "app/Run -> "

	logger.Logger.Process.Println("запуск сервера...")

	logger.Logger.Process.Println("инициализация конфига...")
	if err := initConfig(); err != nil {
		logger.Logger.Error.Fatalln("ошибка инициализации конфига: ", err)
	}
	logger.Logger.Process.Println("инициализация конфига: успешно")

	ctx := context.Background()

	logger.Logger.Process.Println("подключение к серверу nats-streaming...")
	stanConnect, err := stan.Connect(viper.GetString(stanClusterIDKey), viper.GetString(stanClientIDKey))
	if err != nil {
		err = &er.Error{Err: err, Code: er.ErrorNoConnect, Message: "не удалось установить соединение к nats-streaming-server", Op: op}
		logger.Logger.Error.Fatalln(err.Error())
	}
	defer func() {
		if err := stanConnect.Close(); err != nil {
			err = &er.Error{Err: err, Code: er.ErrorClose, Message: "не удалось закрыть соединение с nats-streaming-server", Op: op}
			logger.Logger.Error.Println(err.Error())
		}
	}()
	logger.Logger.Process.Println("подключение к серверу nats-streaming: успешно")

	logger.Logger.Process.Println("подключение к базе данных Postgres...")

	confDB := viper.Get(dbKey).(map[string]interface{})
	url := fmt.Sprintf("postgres://%v:%v@%v:%v/%v",
		confDB["username"], confDB["password"], confDB["host"], confDB["port"], confDB["dbname"],
	)

	postgresDB, err := postgres.NewPostgres(ctx, url)
	if err != nil {
		err = er.AddOp(err, op)
		logger.Logger.Error.Fatalln(err.Error())
	}
	logger.Logger.Process.Println("подключение к базе данных Postgres: успешно")

	logger.Logger.Process.Println("загрузка кеша InMemory...")
	inMemory, err := in_memory.NewInMemory(ctx, postgresDB)
	if err != nil {
		err = er.AddOp(err, op)
		logger.Logger.Error.Fatalln(err.Error())
	}
	logger.Logger.Process.Println("загрузка кеша InMemory: успешно")

	storage := database.NewStorage(postgresDB, inMemory)

	wbService := service.NewWBService(storage)
	logger.Logger.Process.Println("инициализирован сервис Wildberries")

	logger.Logger.Process.Println("попытка подписаться на канал к nats-streaming...")
	stanSubscriber := subscriber.NewSubscriber(&stanConnect, wbService, viper.GetString(subscriberDurableNameKey))
	err = stanSubscriber.Subscribe(channelNameKey)
	if err != nil {
		err = er.AddOp(err, op)
		logger.Logger.Error.Fatalln(err.Error())
	}
	defer func() {
		if err := stanSubscriber.Subscription.Unsubscribe(); err != nil {
			err = &er.Error{Err: err, Code: er.ErrorClose, Message: "не удалось отписаться от канала nats-streaming", Op: op}
			logger.Logger.Error.Println(err.Error())
		}
	}()
	logger.Logger.Process.Println("попытка подписаться на канал к nats-streaming: успешно")

	handler := rest.NewHandler(wbService)
	http.HandleFunc("/", handler.GetOrder)

	logger.Logger.Process.Println("Сервер собран. Сервер приступил к обслуживанию")
	if err := http.ListenAndServe(fmt.Sprintf(":%v", viper.GetString(portKey)), nil); err != nil {
		err = &er.Error{Err: err, Code: er.ErrorListenAndServe, Message: "не удалось обслуживать сервер http", Op: op}
		logger.Logger.Error.Fatalln(err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
