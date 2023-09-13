package nats

import (
	"WBTech_Level0/models"
	"WBTech_Level0/pkg/repository"
	"encoding/json"
	"github.com/nats-io/stan.go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"sync"
)

type NatsConnection struct {
	stan.Conn
	repo *repository.Repository
}

func NewConnection(repo *repository.Repository) (NatsConnection, error) {
	nc, err := stan.Connect(
		viper.GetString("nats.cluster_id"),
		viper.GetString("nats.client_subscriber"),
		stan.NatsURL(
			viper.GetString("nats.url_sub"),
		),
	)
	if err != nil {
		return NatsConnection{}, err
	}
	return NatsConnection{nc, repo}, err
}

func (nc *NatsConnection) Close() {
	if err := nc.Conn.Close(); err != nil {
		logrus.Errorf("[Consumer] error occurred while closing connection to NATS server: %v", err)
	}
}

func (nc *NatsConnection) Subscribe(wg *sync.WaitGroup) error {
	defer wg.Done()
	sc, err := nc.Conn.Subscribe(viper.GetString("nats.subject"), func(msg *stan.Msg) {
		var order models.Order
		logrus.Printf("[NATS connection] Length of received data: %d", len(msg.Data))
		err := json.Unmarshal(msg.Data, &order)
		if err != nil {
			logrus.Errorf("[NATS connection] error occurred while unmarshaling order: %v", err)
			return
		}
		err = nc.repo.CreateOrder(order)
		if err != nil {
			logrus.Errorf("[Consumer] error occurred while posting order to DB: %v", err)
			return
		}
		err = nc.repo.PutOrder(order)
		if err != nil {
			logrus.Errorf("[Consumer] error occurred while putting order data into cache: %v", err)
			return
		}
	})
	if err != nil {
		return err
	}
	for {
		if !sc.IsValid() {
			wg.Done()
			break
		}
	}
	err = sc.Unsubscribe()
	if err != nil {
		return err
	}
	logrus.Println("[Consumer] unsubscribed from the NATS server...")
	return nil
}
