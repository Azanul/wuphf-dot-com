package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Azanul/wuphf-dot-com/notification/internal/controller/notification"
	httphandler "github.com/Azanul/wuphf-dot-com/notification/internal/handler/http"
	"github.com/Azanul/wuphf-dot-com/notification/internal/handler/kafka"
	twiliosms "github.com/Azanul/wuphf-dot-com/notification/internal/integration/twilio-sms"
	"github.com/Azanul/wuphf-dot-com/notification/internal/repository/memory"

	"github.com/IBM/sarama"
)

func main() {
	kafkaBrokers := os.Getenv("KAFKA_BROKERS")

	log.Println("Starting the notification service")
	repo := memory.New()

	twilioIntegration := twiliosms.New()
	ctrl := notification.New(repo)
	ctrl.AddIntegration(twilioIntegration)

	h := httphandler.New(ctrl)

	topics := []string{"chats", "notifications"}

	// Start Kafka consumer
	go func() {
		ctx := context.Background()
		config := sarama.NewConfig()
		config.Consumer.IsolationLevel = sarama.ReadCommitted
		consumerGroup, err := sarama.NewConsumerGroup(strings.Split(kafkaBrokers, ","), "notification_consumer_group", config)
		if err != nil {
			log.Fatalf("Error creating Kafka consumer group: %v", err)
		}
		defer consumerGroup.Close()

		consumer := kafka.New(ctrl)
		err = consumerGroup.Consume(ctx, topics, consumer)
		if err != nil {
			log.Fatalf("Error consuming topic: %v", err)
		}
	}()

	// Endpoints
	http.Handle("/notification", http.HandlerFunc(h.Notification))
	http.Handle("/history", http.HandlerFunc(h.History))

	if err := http.ListenAndServe(":8082", nil); err != nil {
		log.Fatal(err)
	}

	<-make(chan struct{})
}
