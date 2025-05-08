package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/codingbot24-s/db/models"
	
	"github.com/segmentio/kafka-go"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// outbox pattern





func CreateOrder(db *gorm.DB ,order models.Order) error {
	tx := db.Begin()

	if err := tx.Create(&order).Error;err != nil {
		tx.Rollback()
		return err
	}
	
	payload := map[string]interface{}{
		"User_id": order.UserID,
		"Amount": order.Price,
	}
	data,err := json.Marshal(payload)
	if err != nil {
		tx.Rollback()
	}

	event := models.OutBoxEvent{
		EventType: "order_created",
		Payload:   datatypes.JSON(data),
		Processed: false,
		CreatedAt: time.Now(),
	}

	if err := tx.Create(&event).Error; err != nil {
		tx.Rollback()
		return err

	}
	return tx.Commit().Error
}


func StartOutboxProcessor(db *gorm.DB, writer *kafka.Writer) {
	go func() {
		for {
			if err := ProcessOutboxToKafka(db, writer); err != nil {
				log.Printf("Outbox processing failed: %v", err)
			}
			time.Sleep(3 * time.Second) 
		}
	}()
}

func ProcessOutboxToKafka(db *gorm.DB, writer *kafka.Writer) error {
	var events []models.OutBoxEvent
	if err := db.Where("processed = ?", false).Order("created_at").Limit(100).Find(&events).Error; err != nil {
		return fmt.Errorf("failed to fetch outbox events: %w", err)
	}

	for _, event := range events {
		log.Printf("Sending event %s to Kafka...\n", event.ID)

		
		//ka
		retryErr := retry(3, 2*time.Second, func() error {
			return writer.WriteMessages(context.Background(), kafka.Message{
				Key:   []byte(event.EventType),
				Value: event.Payload,
			})
		})

		if retryErr != nil {
			log.Printf("Failed to publish event %s after retries: %v\n", event.ID, retryErr)
			continue 
		}

		if err := db.Model(&event).Update("processed", true).Error; err != nil {
			log.Printf("Failed to mark event %s as processed: %v\n", event.ID, err)
		} else {
			log.Printf("Event %s sent and marked as processed\n", event.ID)
		}
	}

	return nil
}