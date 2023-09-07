package main

import (
	WBTech_Level_0 "WBTech_Level0"
	"WBTech_Level0/configs"
	"WBTech_Level0/pkg/handler"
	"WBTech_Level0/pkg/repository"
	"WBTech_Level0/pkg/service"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"os/signal"
)

func main() {
	err := configs.InitConfig()
	if err != nil {
		logrus.Panicf("error occured initializing configs: %v", err)
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
		logrus.Panicf("error occured connection to DB: %v", err)
	}

	// Закрываем БД
	defer func() {
		if err := db.Close(); err != nil {
			logrus.Panicf("error occured while closing DB: %v", err)
		}
		logrus.Println("closing DB...")
	}()

	repo := repository.NewRepository(db)
	services := service.NewService(repo)
	handlers := handler.NewHandler(services)
	srv := new(WBTech_Level_0.Server)

	//closer.Bind(srv.Shutdown)
	go func() {
		if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
			logrus.Panicf("error occured running consumer server: %v", err)
		}
		logrus.Println("App started...")
	}()

	// TODO проработать нормальную схему закрытия приложения...
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)
	//signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	srv.Shutdown()
	// TODO возможно, стоит вернуть как было, так как в closer вызывается os.Exit -> defer'ы бессмыслены
	//closer.Hold()
	logrus.Println("App closed...")
}
