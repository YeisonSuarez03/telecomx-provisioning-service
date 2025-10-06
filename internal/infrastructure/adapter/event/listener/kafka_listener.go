package listener

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"

	"telecomx-provisioning-service/internal/application/service"
	"telecomx-provisioning-service/internal/domain/model"
)

type CustomerEvent struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type CustomerPayload struct {
	UserID      string `json:"userId"`
	ServiceName string `json:"serviceName"`
	IsActive    bool   `json:"isActive"`
}

func StartKafkaListener(svc *service.ProvisioningService, brokers []string, topic, group, client string) error {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: group,
		Dialer: &kafka.Dialer{
			ClientID: client,
		},
	})
	defer reader.Close()

	log.Printf("[Kafka] Listening on topic: %s", topic)

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Println("Kafka error:", err)
			return err
		}

		// Log raw message metadata for observability
		log.Printf("[Kafka] Received message partition=%d offset=%d key=%q", msg.Partition, msg.Offset, string(msg.Key))

		var event CustomerEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Println("Invalid event:", err)
			continue
		}

		// Log parsed event type and a trimmed payload preview
		log.Printf("[Kafka] Parsed event type=%s payload_len=%d complete=%v", event.Type, len(event.Payload), event)

		var payload CustomerPayload
		_ = json.Unmarshal(event.Payload, &payload)

		switch event.Type {
		case "Customer.Created":
			if payload.ServiceName == "" {
				payload.ServiceName = "Internet"
			}
			err := svc.Create(context.Background(), &model.Provisioning{
				UserID:      payload.UserID,
				ServiceName: payload.ServiceName,
				Status:      "Active",
			})
			if err != nil {
				log.Println("Error creating customer:", err)
				return err
			}
		case "Customer.Updated":
			status := "Active"
			if !payload.IsActive {
				status = "Suspended"
			}
			err := svc.UpdateStatus(context.Background(), payload.UserID, status)
			if err != nil {
				log.Println("Error updating customer status:", err)
				return err
			}
		case "Customer.Deleted":
			err := svc.Delete(context.Background(), payload.UserID)
			if err != nil {
				log.Println("Error deleting customer:", err)
				return err
			}
		}
	}
}
