package main

import (
	"log"

	"github.com/MaximumTroubles/go-todo-app"
	"github.com/MaximumTroubles/go-todo-app/pkg/handler"
	"github.com/MaximumTroubles/go-todo-app/pkg/repository"
	"github.com/MaximumTroubles/go-todo-app/pkg/service"
	"github.com/spf13/viper"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("error intializing configs: %s", err.Error())
	}

	repos := repository.NewRepository()
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(todo.Server)
	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		log.Fatalf("error occurred while running http server: %s", err.Error())
	}

}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
