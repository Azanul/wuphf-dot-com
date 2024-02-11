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
			id, err := c.ctrl.Post(context.TODO(), n["sender"], n["receiver"], n["msg"])
			if err != nil {
				log.Printf("Error creating notification: %v\n", err)
			} else {
				log.Printf("Notification created: %s\n", id)
			}
		} else {
			log.Printf("Error unmarshaling message: %v\n", err)
		}
	}
	return nil
}

// parseMultipartFormData parses the multipart form-data message and returns a map of form values.
func parseMultipartFormData(message []byte) (map[string]string, error) {
	formValues := map[string]string{}
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
		formValues[name] = buf.String()

		part.Close()
	}

	return formValues, nil
}
