package kafka

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log"
	"mime"
	"mime/multipart"

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
		n, err := parseMultipartFormData(msg.Value)
		if err == nil {
			if _, ok := n["chat_id"]; ok {
				id, err := c.ctrl.Post(context.TODO(), n["sender"].(string), n["chat_id"].(string), n["msg"].(string))
				if err != nil {
					log.Printf("Error creating notification: %v\n", err)
				} else {
					log.Printf("Notification created: %s\n", id)
				}
			} else {
				receivers, ok := n["receivers"].([]string)
				if !ok {
					receivers = []string{n["receivers"].(string)}
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

// parseMultipartFormData parses the multipart form-data message and returns a map of form values.
func parseMultipartFormData(message []byte) (map[string]interface{}, error) {
	formValues := map[string]interface{}{}
	reader := multipart.NewReader(bytes.NewReader(message), "--------------------------")
	for {
		part, err := reader.NextPart()
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return nil, err
		}

		var buf bytes.Buffer
		if _, err := io.Copy(&buf, part); err != nil {
			return nil, err
		}

		header := part.Header.Get("Content-Disposition")
		_, params, err := mime.ParseMediaType(header)
		if err != nil {
			return nil, err
		}
		name := params["name"]
		value := buf.String()
		if existingValue, exists := formValues[name]; exists {
			switch v := existingValue.(type) {
			case string:
				formValues[name] = []string{v, value}
			case []string:
				formValues[name] = append(v, value)
			default:
				return nil, errors.New("unsupported type in form values")
			}
		} else {
			formValues[name] = value
		}

		part.Close()
	}

	return formValues, nil
}
