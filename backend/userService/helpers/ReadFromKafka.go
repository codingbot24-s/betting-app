package helpers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
	"gorm.io/gorm"
)


type OrderCreatedEvent struct {
	UserID string `json:"userId"`
	Amount float64 `json:"amount"`
}

func ReadFromKafka(db * gorm.DB) {
	// create a reader
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"localhost:9092"},
		Topic:    "order-created",
		GroupID:  "user-service",
	})
	fmt.Println("Listening for messages on 'order-created' topic")

	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			fmt.Println("error reading message:", err)
			continue
		}

		var event OrderCreatedEvent
		if err := json.Unmarshal(m.Value, &event); err != nil {
			fmt.Println("error unmarshalling message:", err)
			continue
		}	
		
		fmt.Println("Received message:", event)

		// subtract the amount from the user's balance
			// Pass the db
		err = SubtractUserBalance(db , event.UserID, event.Amount)
		if err != nil {
			fmt.Println("error subtracting user balance:", err)
			continue
		}else{
			fmt.Println("User balance updated successfully")
		}
	}
}

