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

	print("good set")

	orders, err := storage.Postgres.GetOrders(ctx)
	if err != nil {
		panic(err)
	}

	for _, order := range *orders {
		fmt.Println(order)
	}
	print("good get")

}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
