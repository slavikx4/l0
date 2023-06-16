package app

import (
	"context"
	"github.com/nats-io/stan.go"
	"github.com/slavikx4/l0/internal/database"
	in_memory "github.com/slavikx4/l0/internal/database/in-memory"
	"github.com/slavikx4/l0/internal/database/postgres"
	"github.com/slavikx4/l0/internal/service"
	"github.com/slavikx4/l0/internal/transport/nats/subscriber"
	er "github.com/slavikx4/l0/pkg/error"
	"github.com/slavikx4/l0/pkg/logger"
	"github.com/spf13/viper"
)

// TODO убрать константы в конфиг
const (
	stanClusterID         = "wbCluster"
	stanClientID          = "wbSubscribe"
	SubscriberDurableName = "name"
	channelName           = "orderChannel"
	portKey               = "port"
	url                   = "postgres://postgres:postgres@localhost:5432/L0"
)

// TODO обернуть инициализацию всех структур в функции
func Run() {
	const op = "app/Run -> "
	//logger.Logger.Process.Println("запуск сервера")
	//
	//if err := initConfig(); err != nil {
	//	logger.Logger.Error.Fatalln("ошибка инициализации конфига: ", err)
	//}
	//
	//ctx := context.Background()
	//
	//postgres, err := postgres.NewPostgres(ctx, "postgres://postgres:postgres@localhost:5432/L0")
	//if err != nil {
	//	panic(err)
	//}
	//storage := database.Storage{
	//	Postgres: postgres,
	//	InMemory: in_memory.NewInMemory(),
	//}
	//data, err := os.ReadFile("model.json")
	//if err != nil {
	//	panic(err)
	//}
	//order := models.Order{}
	//if err := json.Unmarshal(data, &order); err != nil {
	//	panic(err)
	//}
	//
	//if err := storage.Postgres.SetOrder(ctx, &order); err != nil {
	//	panic(err)
	//}
	//
	//fmt.Println("good set in postgres")
	//
	//orders, err := storage.Postgres.GetOrders(ctx)
	//if err != nil {
	//	panic(err)
	//}
	//
	//fmt.Println("good get from postgres")
	//
	//storage.InMemory.SetOrder(orders)
	//fmt.Println("good set into in-memory")
	//
	//orderInto, err := storage.InMemory.GetOrder("b563feb7b2b84b6test")
	//if err != nil {
	//	fmt.Println(err.Error())
	//}
	//fmt.Println("good get from in-memory")
	//fmt.Println(orderInto)

	ctx := context.Background()

	stanConnect, err := stan.Connect(stanClusterID, stanClientID)
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

	postgresDB, err := postgres.NewPostgres(ctx, url)
	if err != nil {
		err = er.AddOp(err, op)
		logger.Logger.Error.Fatalln(err.Error())
	}

	inMemory, err := in_memory.NewInMemory(ctx, postgresDB)
	if err != nil {
		err = er.AddOp(err, op)
		logger.Logger.Error.Fatalln(err.Error())
	}

	storage := database.NewStorage(postgresDB, inMemory)

	wbService := service.NewWBService(storage)

	stanSubscriber := subscriber.NewSubscriber(&stanConnect, wbService, SubscriberDurableName)
	stanSubscriber.Subscribe(channelName)
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
