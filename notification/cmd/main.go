package main

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
	"wuphf.com/notification/internal/controller/notification"
	httphandler "wuphf.com/notification/internal/handler/http"
	"wuphf.com/notification/internal/repository/memory"
)

func main() {
	log.Println("Starting the notification service")
	repo := memory.New()
	ctrl := notification.New(repo)
	h := httphandler.New(ctrl)

	// Start Kafka consumer
	go func() {
		config := sarama.NewConfig()
		consumer, err := sarama.NewConsumer([]string{"localhost:9092"}, config)
		if err != nil {
			log.Fatalf("Error creating Kafka consumer: %v", err)
		}
		defer consumer.Close()

		topic := "notifications"
		partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
		if err != nil {
			log.Fatalf("Error consuming partition: %v", err)
		}
		defer partitionConsumer.Close()

		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	ConsumerLoop:
		for {
			select {
			case msg := <-partitionConsumer.Messages():
				n, err := parseMultipartFormData(msg.Value)
				if err == nil {
					id, err := ctrl.Post(context.TODO(), n["sender"], n["receiver"], n["msg"])
					if err != nil {
						log.Printf("Error creating notification: %v\n", err)
					} else {
						log.Printf("Notification created: %s\n", id)
					}
				} else {
					log.Printf("Error unmarshaling message: %v\n", err)
				}

			case err := <-partitionConsumer.Errors():
				log.Printf("Error consuming message: %v\n", err)
			case <-signals:
				break ConsumerLoop
			}
		}
	}()

	// Endpoints
	http.Handle("/notification", http.HandlerFunc(h.Notification))
	http.Handle("/history", http.HandlerFunc(h.History))

	if err := http.ListenAndServe(":8082", nil); err != nil {
		panic(err)
	}

	<-make(chan struct{})
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
