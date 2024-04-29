package publisher

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"time"
	"wbTechL0/internal/config"
	"wbTechL0/models"

	"github.com/nats-io/nats.go"
)

func PublishOrders(nc *nats.Conn) {
	orders, err := getOrders()
	if err != nil {
		log.Println(err)
		return
	}

	config, err := config.LoadConfig("config/config.json")
	if err != nil {
		log.Fatal("Error loading configuration:", err)
	}

	for _, oneOrder := range orders {

		// create random message intervals to slow down
		r := rand.Intn(1500)
		time.Sleep(time.Duration(r) * time.Millisecond)

		orderString, err := json.Marshal(oneOrder)
		if err != nil {
			log.Println(err)
			continue
		}

		err = nc.Publish(config.Nats.StreamName, []byte(orderString))
		if err != nil {
			log.Println(err)
		} else {
			log.Println("Publisher send message")
		}
	}
}

func getOrders() ([]models.Order, error) {
	rawOrder, _ := ioutil.ReadFile("./orders.json")
	var orderObj []models.Order
	err := json.Unmarshal(rawOrder, &orderObj)

	return orderObj, err
}
