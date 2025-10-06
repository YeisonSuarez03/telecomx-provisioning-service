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

// AltCustomerEvent supports incoming messages that use keys: "event" and "data"
type AltCustomerEvent struct {
    Event string          `json:"event"`
    Data  json.RawMessage `json:"data"`
}

// toJSON returns a compact JSON representation of any value for logging.
// Falls back to an error string if marshaling fails.
func toJSON(v interface{}) string {
    b, err := json.Marshal(v)
    if err != nil {
        return "<json_error: " + err.Error() + ">"
    }
    return string(b)
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
        // Print the raw value as UTF-8 text (json.Marshal on []byte shows base64)
        log.Printf("[Kafka] Raw Value bytes as string=%s", string(msg.Value))
        log.Printf("[Kafka] MESSAGE=%s", toJSON(msg))


        var event CustomerEvent
        if err := json.Unmarshal(msg.Value, &event); err != nil {
            // Try fallback shape {"event": "...", "data": {...}}
            var alt AltCustomerEvent
            if err2 := json.Unmarshal(msg.Value, &alt); err2 != nil {
                log.Println("Invalid event (both shapes failed):", err, "|", err2)
                continue
            }
            event.Type = alt.Event
            event.Payload = alt.Data
        } else if event.Type == "" && len(event.Payload) == 0 {
            // Some producers may use different keys; attempt fallback even if first unmarshal succeeded but empty
            var alt AltCustomerEvent
            if err2 := json.Unmarshal(msg.Value, &alt); err2 == nil {
                event.Type = alt.Event
                event.Payload = alt.Data
            }
        }

        // Log parsed event type and a trimmed payload preview
        log.Printf("[Kafka] Parsed event type=%s payload_len=%d", event.Type, len(event.Payload))
        // Log full event JSON for debugging
        log.Printf("[Kafka] Event JSON=%s", toJSON(event))

		var payload CustomerPayload
		_ = json.Unmarshal(event.Payload, &payload)
        // Log full payload JSON for debugging
        log.Printf("[Kafka] Payload JSON=%s", toJSON(payload))

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
