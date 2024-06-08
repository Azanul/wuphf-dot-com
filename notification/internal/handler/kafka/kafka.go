package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Azanul/wuphf-dot-com/notification/internal/controller/notification"

	"github.com/IBM/sarama"
)

// Handler defines a notification Kafka message handler
type Handler struct {
	ctrl *notification.Controller
}

// New creates a new notification Kafka message handler
func New(ctrl *notification.Controller) *Handler {
	return &Handler{ctrl}
}

func (c Handler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (c Handler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (c Handler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		sess.MarkMessage(msg, "")
		sess.Commit()
		n, err := parseJSON(msg.Value)
		if err == nil {
			if _, ok := n["chat_id"]; ok {
				id, err := c.ctrl.Post(context.TODO(), n["sender"].(string), n["chat_id"].(string), n["msg"].(string))
				if err != nil {
					log.Printf("Error creating notification: %v\n", err)
				} else {
					log.Printf("Notification created: %s\n", id)
				}
			} else {
				receivers, err := parseReceivers(n["receivers"])
				if err != nil {
					log.Printf("Error parsing receivers: %v\n", err)
					continue
				}

				id := c.ctrl.PostChat(context.TODO(), n["sender"].(string), receivers)
				log.Printf("Chat created: %s\n", id)
			}
		} else {
			log.Printf("Error unmarshaling message: %v\n", err)
		}
	}
	return nil
}

// parseJSON parses the JSON message and returns a map of values.
func parseJSON(message []byte) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := json.Unmarshal(message, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func parseReceivers(receivers interface{}) ([]string, error) {
	var result []string
	switch v := receivers.(type) {
	case []interface{}:
		for _, r := range v {
			if str, ok := r.(string); ok {
				result = append(result, str)
			} else {
				return nil, fmt.Errorf("invalid receiver type: %v", r)
			}
		}
	case string:
		result = []string{v}
	default:
		return nil, fmt.Errorf("unexpected type for receivers: %T", v)
	}
	return result, nil
}
