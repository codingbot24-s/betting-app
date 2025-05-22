package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/codingbot24-s/db/models"
	"github.com/segmentio/kafka-go"
	"gorm.io/gorm"
)

type KafkaEvent struct {
	MarketID string
	Outcome  bool
}

func ProcessMarketOutcome(marketID string, outcome bool, db *gorm.DB) error {
	var positions []models.Position

	// Fetch all positions for this market
	result := db.Where("market_id = ?", marketID).Find(&positions)
	if result.Error != nil {
		return fmt.Errorf("error fetching positions: %v", result.Error)
	}

	fmt.Printf("Found %d positions for market %s\n", len(positions), marketID)

	// Process winners and losers
	for _, position := range positions {
		isWinner := (outcome && position.Side == "YES") || (!outcome && position.Side == "NO")

		if isWinner {
			// Calculate winnings (example: winners get double their bet)
			winnings := position.Amount * 2

			fmt.Printf("Winner found: UserID: %s, Amount: %.2f, Winnings: %.2f\n",
				position.UserID, position.Amount, winnings)

			// TODO: Call balance service to credit winnings
			err := UpdateUserBalance(position.UserID, winnings)
			if err != nil {
				fmt.Printf("Error updating balance for user %s: %v\n", position.UserID, err)
				continue
			}
		}
	}

	return nil
}

func UpdateUserBalance(userID string, amount float64) error {
	// TODO: Implement balance update logic through balance service
	return nil
}

func ReadFromKafka(db *gorm.DB) {
	// Configure Kafka reader with better error handling
	config := kafka.ReaderConfig{
		Brokers:     []string{"localhost:9092"},
		Topic:       "market-resolved",
		GroupID:   "trading-service-consumer-group",
	
		
		// Add error handling
		ErrorLogger: kafka.LoggerFunc(func(s string, i ...interface{}) {
			fmt.Printf("Kafka error: "+s+"\n", i...)
		}),
		// Add info logging
		Logger: kafka.LoggerFunc(func(s string, i ...interface{}) {
			fmt.Printf("Kafka info: "+s+"\n", i...)
		}),
	}

	reader := kafka.NewReader(config)
	defer reader.Close()

	fmt.Println("Starting to read from Kafka topic: market-resolved")

	ctx := context.Background()
	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			fmt.Printf("Error reading message: %v\n", err)
			// Add delay before retrying
			time.Sleep(time.Second * 5)
			continue
		}

		var event KafkaEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			fmt.Printf("Error unmarshalling message: %v\n", err)
			continue
		}

		// validate the event
		if event.MarketID == "" {
			fmt.Printf("Invalid market ID: %s\n", event.MarketID)
			continue
		}

		fmt.Printf("Processing event - Market: %s, Outcome: %v\n", event.MarketID, event.Outcome)

		if err := ProcessMarketOutcome(event.MarketID, event.Outcome, db); err != nil {
			fmt.Printf("Error processing market outcome: %v\n", err)
			continue
		}

		// Commit the message after successful processing
		if err := reader.CommitMessages(ctx, msg); err != nil {
			fmt.Printf("Error committing message: %v\n", err)
		}
	}
}
