package main

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/Azanul/wuphf-dot-com/user/gen"
	"github.com/Azanul/wuphf-dot-com/user/pkg/model"

	"github.com/IBM/sarama"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type customString string

var userString customString = "user"

var (
	ErrNoMetadata   = errors.New("no metadata found in context")
	ErrInvalidToken = errors.New("invalid token")
)

// Route represents a route configuration
type Route struct {
	Path       string
	BackendURL string
	Handler    RouteHandler
	match      func(string) bool
}

type RouteHandler interface {
	ServeHTTP(rw http.ResponseWriter, req *http.Request)
}

// KafkaMessageProducer represents a Kafka message producer handler
type KafkaMessageProducer struct {
	KafkaTopic string
	Producer   sarama.AsyncProducer
}

func (kafkaProducer *KafkaMessageProducer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	message := &sarama.ProducerMessage{
		Topic: kafkaProducer.KafkaTopic,
		Value: sarama.StringEncoder(body),
	}

	kafkaProducer.Producer.BeginTxn()
	kafkaProducer.Producer.Input() <- message

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Message produced successfully"))
}

// Gateway represents the API gateway
type Gateway struct {
	Routes          []*Route
	SecureRoutes    []*Route
	AuthServiceAddr string
}

// NewGateway initializes a new API gateway
func NewGateway(authServiceAddr string) *Gateway {
	return &Gateway{
		Routes:          []*Route{},
		SecureRoutes:    []*Route{},
		AuthServiceAddr: authServiceAddr,
	}
}

// AddRoute adds a route to the gateway
func (gateway *Gateway) AddRoute(path, backendURL string, cmp func(string, string) bool, secure bool, handler RouteHandler) {
	route := &Route{
		Path:       path,
		BackendURL: backendURL,
		Handler:    handler,
		match: func(url string) bool {
			return cmp(url, path)
		},
	}
	if secure {
		gateway.SecureRoutes = append(gateway.SecureRoutes, route)
	} else {
		gateway.Routes = append(gateway.Routes, route)
	}
}

// ServeHTTP handles incoming HTTP requests
func (gateway *Gateway) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range gateway.Routes {
		if route.match(r.URL.Path) {
			route.Handler.ServeHTTP(w, r)
			return
		}
	}
	for _, route := range gateway.SecureRoutes {
		if route.match(r.URL.Path) {
			user, err := gateway.authenticate(r)
			if err != nil {
				log.Println("Error authenticating:", err)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), userString, user)
			r = r.WithContext(ctx)
			route.Handler.ServeHTTP(w, r)
			return
		}
	}
	http.NotFound(w, r)
}

// authenticate performs authentication using gRPC metadata
func (gateway *Gateway) authenticate(r *http.Request) (*model.User, error) {
	token := r.Header.Get("Authorization")
	return gateway.ValidateToken(r.Context(), token)
}

// ValidateToken calls the authentication service to validate the token
func (gateway *Gateway) ValidateToken(ctx context.Context, token string) (*model.User, error) {
	conn, err := grpc.Dial(
		gateway.AuthServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := gen.NewAuthServiceClient(conn)

	const maxRetries = 5
	for i := 0; i < maxRetries; i++ {
		resp, err := client.ValidateToken(ctx, &gen.TokenRequest{Token: token})
		if err != nil {
			if shouldRetry(err) {
				continue
			}
			return nil, err
		}
		if resp.GetValid() {
			return model.UserFromProto(resp.GetUser()), nil
		}
		return nil, ErrInvalidToken
	}
	return nil, errors.New("maximum retry attempts reached")
}

// shouldRetry checks if the error is retryable
func shouldRetry(err error) bool {
	e, ok := status.FromError(err)
	if !ok {
		return false
	}
	return e.Code() == codes.DeadlineExceeded || e.Code() == codes.ResourceExhausted || e.Code() == codes.Unavailable
}

func main() {
	userService := "http://localhost:8081"
	notificationService := "http://localhost:8082"
	authService := "localhost:50051"
	gateway := NewGateway(authService)

	// Kafka producer setup
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.RequiredAcks = sarama.WaitForAll         // Wait for all replicas to acknowledge the record
	kafkaConfig.Producer.Compression = sarama.CompressionSnappy   // Compress messages
	kafkaConfig.Producer.Flush.Frequency = 100 * time.Millisecond // Flush batches every 100ms
	kafkaConfig.Producer.Idempotent = true                        // Idempotent producer
	kafkaConfig.Net.MaxOpenRequests = 1                           // Only one outstanding request
	kafkaProducer, err := sarama.NewAsyncProducer([]string{"kkafka:9092"}, kafkaConfig)
	if err != nil {
		log.Fatalf("Failed to start Kafka producer: %v", err)
	}
	log.Println("Kafka async producer started")
	defer func() {
		if err := kafkaProducer.Close(); err != nil {
			log.Fatalf("Error closing Kafka producer: %v", err)
		}
	}()

	// Routes
	gateway.AddRoute("/user", userService, strings.EqualFold, false, httputil.NewSingleHostReverseProxy(MustParse(userService)))
	gateway.AddRoute("/auth", userService, strings.HasPrefix, false, httputil.NewSingleHostReverseProxy(MustParse(userService)))
	gateway.AddRoute("/notification", notificationService, strings.EqualFold, true, &KafkaMessageProducer{KafkaTopic: "notifications", Producer: kafkaProducer})

	http.Handle("/", gateway)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func MustParse(backendURL string) *url.URL {
	backend, err := url.Parse(backendURL)
	if err != nil {
		log.Fatalf("Failed to parse backend URL: %s", err)
	}
	return backend
}
