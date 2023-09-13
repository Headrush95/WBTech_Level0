package main

import (
	"github.com/nats-io/stan.go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io"
	"net/http"
)

const port = ":3000"
const natsURL = "nats://localhost:4222"

/*
Publisher нужен только для того, чтобы проверить работоспособность брокера
для передачи информации о заказе.
*/
type natsConfig struct {
	ClusterId string
	UrlSub    string
	UrlPub    string
	ClientSub string
	ClientPub string
	Subject   string
}

func configInit() natsConfig {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	logrus.SetFormatter(new(logrus.JSONFormatter))
	if err != nil {
		logrus.Panicf("[Publisher] error occurred during reading config file: %v", err)
	}
	return natsConfig{
		ClusterId: viper.GetString("nats.cluster_id"),
		UrlSub:    viper.GetString("nats.url_sub"),
		UrlPub:    viper.GetString("nats.url_pub"),
		ClientSub: viper.GetString("nats.client_subscriber"),
		ClientPub: viper.GetString("nats.client_producer"),
		Subject:   viper.GetString("nats.subject"),
	}
}

func main() {
	http.HandleFunc("/store", storeOrder)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		panic(err)
	}
	logrus.Println("[Publisher] successfully started")
}

func storeOrder(w http.ResponseWriter, r *http.Request) {
	config := configInit()

	natsConn, err := stan.Connect(
		config.ClusterId,
		config.ClientPub,
		stan.NatsURL(config.UrlPub),
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.Errorf("[Publisher] error occurred during connecting to NATS server: %v", err)
		return
	}
	defer func(nc *stan.Conn) {
		if err := (*nc).Close(); err != nil {
			logrus.Errorf("[Publisher] error occurred during closing connection to NATS server: %v", err)
		}
	}(&natsConn)

	logrus.Println("[Publisher] successfully connected to NATS server")

	//var order []byte
	//_, err = r.Body.Read(order)
	order, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.Panicf("[Publisher] error occurred during reading request body: %v", err)
	}
	defer r.Body.Close()
	logrus.Printf("[Publisher] successfully read request body, length of data = %d", len(order))
	err = natsConn.Publish(config.Subject, order)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.Errorf("[Publisher] error occurred during publishing to NATS server: %v", err)
		return
	}
	logrus.Println("[Publisher] successfully sent order info")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Order info has been sent"))
}
