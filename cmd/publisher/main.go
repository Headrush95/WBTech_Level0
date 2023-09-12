package main

import (
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
	"net/http"
)

const port = ":3000"
const natsURL = "nats://localhost:4222"

/*
Publisher нужен только для того, чтобы проверить работоспособность брокера
для передачи информации о заказе.
*/
func main() {
	http.HandleFunc("/store", storeOrder)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		panic(err)
	}

}

func storeOrder(w http.ResponseWriter, r *http.Request) {
	natsConn, err := nats.Connect(natsURL)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.Panicf("[Publisher] error occurred during connecting to NATS server: %v", err)
	}
	defer natsConn.Close()
	logrus.Println("[Publisher] successfully connected to NATS server")

	var order []byte
	_, err = r.Body.Read(order)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.Panicf("[Publisher] error occurred during reading request body: %v", err)
	}
	defer r.Body.Close()

	err = natsConn.Publish("orders", order)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.Panicf("[Publisher] error occurred during publishing to NATS server: %v", err)
	}
	logrus.Println("[Publisher] successfully sent order info")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Order info has been sent"))
}
