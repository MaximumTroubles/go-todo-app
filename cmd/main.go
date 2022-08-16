package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/MaximumTroubles/go-todo-app"
	"github.com/MaximumTroubles/go-todo-app/pkg/handler"
	"github.com/MaximumTroubles/go-todo-app/pkg/repository"
	"github.com/MaximumTroubles/go-todo-app/pkg/service"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	//here we make Json format for logging that more preferred for log collectors system like (Graylog, Logstash etc.)
	logrus.SetFormatter(new(logrus.JSONFormatter))

	// viper package init it's read config file that we created and take data from it
	if err := initConfig(); err != nil {
		logrus.Fatalf("error intializing configs: %s", err.Error())
	}

	//.env environment load
	if err := godotenv.Load(".env"); err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())
	}

	// fulfill Postgres config here and through the  sqlx driver make connection to the postgres container
	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBname:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})

	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}

	// init new instance of repository that take as argument pointer to sqlx.DB struct which we define and get from Connect
	// method of repository.NewPostgresDb(). And return pointer of Repository struct with concrete entities repositories as pointer as well.
	repos := repository.NewRepository(db)
	services := service.NewService(repos)

	// Here we are creating a new handler and as argument provide Service Struct with Interfaces
	handlers := handler.NewHandler(services)

	srv := new(todo.Server)
	go func() {

		// method Run() blocking work of main's goroutine because method http.ListenAndServe launched endless for() loop for
		// incoming http requests

		// But now we launched server on go routine that's not blocking execution of function main() and we just quit  app
		if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
			logrus.Fatalf("error occurred while running http server: %s", err.Error())
		}
	}()

	logrus.Print("TodoApp Started")

	// to  avoid quiting an app implement blocking main() function with os.Signal channel help
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("TodoApp Finished")
	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured on server shutting down: %s", err.Error())
	}
	if err := db.Close(); err != nil {
		logrus.Errorf("error occured on db connection close: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
