package main

import (
	"context"
	"log"
	"net/http"

	"github.com/Azanul/wuphf-dot-com/notification/internal/controller/notification"
	httphandler "github.com/Azanul/wuphf-dot-com/notification/internal/handler/http"
	"github.com/Azanul/wuphf-dot-com/notification/internal/handler/kafka"
	"github.com/Azanul/wuphf-dot-com/notification/internal/repository/memory"

	"github.com/IBM/sarama"
)

func main() {
	log.Println("Starting the notification service")
	repo := memory.New()
	ctrl := notification.New(repo)
	h := httphandler.New(ctrl)

	// Start Kafka consumer
	go func() {
		ctx := context.Background()
		config := sarama.NewConfig()
		config.Consumer.IsolationLevel = sarama.ReadCommitted
		consumerGroup, err := sarama.NewConsumerGroup([]string{"kkafka:9092"}, "notification_consumer_group", config)
		if err != nil {
			log.Fatalf("Error creating Kafka consumer group: %v", err)
		}
		defer consumerGroup.Close()

		topic := "notifications"
		consumer := kafka.New(ctrl)
		err = consumerGroup.Consume(ctx, []string{topic}, consumer)
		if err != nil {
			log.Fatalf("Error consuming topic: %v", err)
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
