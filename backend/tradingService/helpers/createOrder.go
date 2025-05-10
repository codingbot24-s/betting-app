package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/codingbot24-s/db/models"
	"github.com/google/uuid"

	"github.com/segmentio/kafka-go"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// outbox pattern

func CreateOrder(db *gorm.DB, order models.Order) error {

	tx := db.Begin()

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		return err
	}

	// need to calculate total

	order.Total = order.Price * float64(order.Quantity)

	payload := map[string]interface{}{
		"User_id": order.UserID,
		"Amount":  order.Total,
	}

	fmt.Println("payload", payload)
	data, err := json.Marshal(payload)
	if err != nil {
		tx.Rollback()
	}

	event := models.OutBoxEvent{
		ID:        uuid.New(),
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

		//kafka with retry
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

		// update the outbox set processed to true
		if err := db.Model(&models.OutBoxEvent{}).
			Where("id = ?", event.ID).
			Update("processed", true).Error; err != nil {
			log.Printf("Failed to mark event %s as processed: %v\n", event.ID, err)
		} else {
			log.Printf("Event %s sent and marked as processed\n", event.ID)
		}
	}

	return nil
}
