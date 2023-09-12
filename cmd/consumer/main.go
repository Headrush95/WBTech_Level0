package main

import (
	WBTech_Level_0 "WBTech_Level0"
	"WBTech_Level0/configs"
	"WBTech_Level0/pkg/handler"
	"WBTech_Level0/pkg/repository"
	"WBTech_Level0/pkg/service"
	"context"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	err := configs.InitConfig()
	if err != nil {
		logrus.Panicf("[Consumer] error occured initializing configs: %v", err)
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		logrus.Panicf("[Consumer] error occured connection to DB: %v", err)
	}

	// Закрываем БД
	defer func() {
		if err := db.Close(); err != nil {
			logrus.Panicf("[Consumer] error occured while closing DB: %v\n", err)
		}

		logrus.Println("[Consumer] closing DB...")
	}()

	repo := repository.NewRepository(db)
	defer func() {
		err := repo.PostgresRepository.CloseStatements()
		if err != nil {
			logrus.Printf("[Consumer] error occurred while closing statements: %v\n", err)
		}
	}()

	services := service.NewService(repo)
	handlers := handler.NewHandler(services)
	srv := new(WBTech_Level_0.Server)

	// Запускаем http сервер
	go func() {
		if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
			logrus.Panicf("[Consumer] error occured running consumer server: %v\n", err)
		}
		logrus.Println("[Consumer] App started...")
	}()

	// TODO проработать нормальную схему закрытия приложения...
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()
	stop()
	logrus.Println("[Consumer] shutting down gracefully, press Ctrl+C again to force")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = srv.Shutdown(ctx); err != nil {
		logrus.Panicf("[Consumer] server forced to shutdown: %v\n", err)
	}

	logrus.Println("[Consumer] server exiting...")
}
