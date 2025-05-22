package helpers

import (
	"context"
	"encoding/json"
	"fmt"


	"github.com/segmentio/kafka-go"
	"gorm.io/gorm"
)

type OrderCreatedEvent struct {
	UserID string  `json:"User_id"` // Changed to match JSON field "User_id"
	Amount float64 `json:"Amount"`  // Changed to match JSON field "Amount"
}

func ReadFromKafka(db *gorm.DB) {
	// create a reader
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "order.created",
		GroupID: "user-service",
	})
	fmt.Println("Listening for messages on 'order.created' topic")

	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			fmt.Printf("Error reading message: %v\n", err)
			continue
		}

		// Log raw message
		fmt.Printf("Raw message received: %s\n", string(m.Value))

		var event OrderCreatedEvent
		if err := json.Unmarshal(m.Value, &event); err != nil {
			fmt.Printf("Error unmarshalling message: %v\n", err)
			continue
		}

		// Validate UserID is not empty
		if event.UserID == "" {
			fmt.Println("Warning: Received message with empty UserID")
			continue
		}

	
		//
		fmt.Println("subtarcting amount from user balance", event.UserID) 		
		balance, err := SubtractUserBalanceAndReturnCurrentBalance(db, event.UserID, event.Amount)
		if err != nil {
			fmt.Printf("Error subtracting balance for user %s: %v\n", event.UserID, err)
		}
		fmt.Println("balance of user updated successfully", event.UserID)

		// current blance 
		fmt.Printf("Current balance for user %s: %.2f\n", event.UserID, balance)
		
	}
}
