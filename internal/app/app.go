package app

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/slavikx4/l0/internal/database"
	in_memory "github.com/slavikx4/l0/internal/database/in-memory"
	"github.com/slavikx4/l0/internal/database/postgres"
	"github.com/slavikx4/l0/internal/models"
	"github.com/spf13/viper"
	"os"
)

const (
	portKey = "port"
)

func Run() {
	//logger.Logger.Process.Println("запуск сервера")
	//
	//if err := initConfig(); err != nil {
	//	logger.Logger.Error.Fatalln("ошибка инициализации конфига: ", err)
	//}

	ctx := context.Background()

	postgres, err := postgres.NewPostgres(ctx, "postgres://postgres:postgres@localhost:5432/L0")
	if err != nil {
		panic(err)
	}
	storage := database.Storage{
		Postgres: postgres,
		InMemory: in_memory.NewInMemory(),
	}
	data, err := os.ReadFile("model.json")
	if err != nil {
		panic(err)
	}
	order := models.Order{}
	if err := json.Unmarshal(data, &order); err != nil {
		panic(err)
	}

	if err := storage.Postgres.SetOrder(ctx, &order); err != nil {
		panic(err)
	}

	fmt.Println("good set in postgres")

	orders, err := storage.Postgres.GetOrders(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println("good get from postgres")

	storage.InMemory.SetOrder(orders)
	fmt.Println("good set into in-memory")

	orderInto, err := storage.InMemory.GetOrder("b563feb7b2b84b6test")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("good get from in-memory")
	fmt.Println(orderInto)

}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
