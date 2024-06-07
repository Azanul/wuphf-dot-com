package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Azanul/wuphf-dot-com/user/gen"
	"github.com/Azanul/wuphf-dot-com/user/internal/controller/user"
	grpchandler "github.com/Azanul/wuphf-dot-com/user/internal/handler/grpc"
	httphandler "github.com/Azanul/wuphf-dot-com/user/internal/handler/http"
	"github.com/IBM/sarama"

	"google.golang.org/grpc"

	"github.com/Azanul/wuphf-dot-com/user/internal/repository/memory"
)

func main() {
	log.Println("Starting the user service")
	kafkaBrokers := os.Getenv("KAFKA_BROKERS")

	// Kafka producer setup
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.RequiredAcks = sarama.WaitForAll         // Wait for all replicas to acknowledge the record
	kafkaConfig.Producer.Compression = sarama.CompressionSnappy   // Compress messages
	kafkaConfig.Producer.Flush.Frequency = 100 * time.Millisecond // Flush batches every 100ms
	kafkaConfig.Producer.Idempotent = true                        // Idempotent producer
	kafkaConfig.Net.MaxOpenRequests = 1                           // Only one outstanding request
	kafkaProducer, err := sarama.NewAsyncProducer(strings.Split(kafkaBrokers, ","), kafkaConfig)
	if err != nil {
		log.Fatalf("Failed to start Kafka producer: %v", err)
	}
	log.Println("Kafka async producer started")
	defer func() {
		if err := kafkaProducer.Close(); err != nil {
			log.Fatalf("Error closing Kafka producer: %v", err)
		}
	}()

	repo := memory.New()
	ctrl := user.New(repo, kafkaProducer)

	h := httphandler.New(ctrl)
	g := grpchandler.New(ctrl)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 50051))
	if err != nil {
		log.Fatal("Failed to listen:", err)
	}
	srv := grpc.NewServer()
	gen.RegisterAuthServiceServer(srv, g)
	go func() {
		if err := srv.Serve(lis); err != nil {
			log.Fatal("Failed to start the gRPC server:", err)
		}
	}()

	// Endpoints
	http.Handle("/user", http.HandlerFunc(h.User))
	http.Handle("/auth/register", http.HandlerFunc(h.Register))
	http.Handle("/auth/login", http.HandlerFunc(h.Login))

	if err := http.ListenAndServe(":8081", nil); err != nil {
		panic(err)
	}
}
