package helpers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
	"gorm.io/gorm"
)

func StartProcessor(ctx context.Context, db *gorm.DB, writer *kafka.Writer) {
	go func() {
		ticker := time.NewTicker(3 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := ProccessOutboxToKafka(db, writer); err != nil {
					log.Printf("Error processing outbox: %v", err)
				}
			}
		}
	}()
}

// TODO: write a reader for kafka in trading service for reading the market resolved event it will give a marketid and outcome fetch all the positions for that market id and if outcome is true then winner are those who have chose the yes
type OutboxEvent struct {
	ID        string `gorm:"primaryKey"`
	EventType string
	Payload   string
	Processed bool
	CreatedAt time.Time
}

func ProccessOutboxToKafka(db *gorm.DB, writer *kafka.Writer) error {
	var events []OutboxEvent
	err := db.Where("processed = ?", false).Order("created_at").Find(&events).Error
	if err != nil {
		return fmt.Errorf("failed to fetch unprocessed events: %w", err)
	}

	for _, event := range events {
		if err := processEvent(db, writer, event); err != nil {
			log.Printf("Error processing event %s: %v", event.ID, err)
			continue
		}
	}
	return nil
}

func processEvent(db *gorm.DB, writer *kafka.Writer, event OutboxEvent) error {
	err := writer.WriteMessages(context.Background(), kafka.Message{

		Key:   []byte(event.ID),
		Value: []byte(event.Payload),
	})

	if err != nil {
		return fmt.Errorf("failed to write message to kafka: %w", err)
	}

	err = db.Model(&OutboxEvent{}).Where("id = ?", event.ID).Update("processed", true).Error
	if err != nil {
		return fmt.Errorf("failed to update propositions for that market id and if outcome is true then winner are those who have chose the yescessed status: %w", err)
	}

	return nil
}
