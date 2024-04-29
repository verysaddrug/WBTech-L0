package consumer

import (
	"encoding/json"
	"log"

	"wbTechL0/internal/storage/pgsql"
	"wbTechL0/models"

	"github.com/nats-io/nats.go"
)

func consumeOrder(jsonOrder []byte, db *pgsql.Database) error {
	var order models.Order

	err := json.Unmarshal(jsonOrder, &order)
	if err != nil {
		log.Println("Error unmarshalling JSON:", err)
		return err
	}

	err = db.SaveOrder(order)
	if err != nil {
		log.Println("Error saving order to the database:", err)
		return err
	}

	log.Println("Order saved successfully.")

	return nil
}

func SubscribeAndConsume(db *pgsql.Database, nc *nats.Conn, streamName string) (*nats.Subscription, error) {
	sub, err := nc.Subscribe(streamName, func(msg *nats.Msg) {
		err := consumeOrder(msg.Data, db)
		if err != nil {
			log.Println("Error consuming order:", err)
		}
	})
	if err != nil {
		log.Fatal("Error subscribing to NATS stream:", err)
		return nil, err
	}

	return sub, nil
}
