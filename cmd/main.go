package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	maksimzhashkevychtodoapp "github.com/LionJr/todo-app"
	"github.com/LionJr/todo-app/pkg/handler"
	"github.com/LionJr/todo-app/pkg/repository"
	"github.com/LionJr/todo-app/pkg/service"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializinng configs: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())
	}

	db, err := repository.NewPostgresDb(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})

	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}

	// gaytaryan znaceniyasy *gin.Engine-a degisli
	// ServeHTTP(ResponseWriter, *Request) metody bolany ucin
	// ol http.Handler interface-y realizowat edyar
	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(maksimzhashkevychtodoapp.Server)
	go func() {
		if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
			logrus.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()

	logrus.Print("TodoApp Started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("TodoApp Shutting Down")

	// server ocuryar
	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured on server shutting down: %s", err.Error())
	}

	// db connection kesyar
	if err := db.Close(); err != nil {
		logrus.Errorf("error occured on db connection close: %s", err.Error())
	}
}

func initConfig() error {
	// config file-lar bilen islemek ucin kop ulanylyan library
	// ilki bilen go get -u github.com/spf13/viper download etmeli
	// bu yerde fayl haysy papkada yerlesyan bolsa son ady yazylyar

	viper.AddConfigPath("configs")
	// AddConfigPath() papkan adyny gorkezenimizden sonra son icinde yerlesyan file adyny bermeli (.extenshion gerek dal) argument edip
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
