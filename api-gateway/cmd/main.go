package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"wuphf.com/user/gen"
	"wuphf.com/user/pkg/model"
)

type customString string

var userString customString = "user"

var (
	ErrNoMetadata   = errors.New("no metadata found in context")
	ErrUnauthorized = errors.New("unauthorized")
)

// Route represents a route configuration
type Route struct {
	Path         string
	BackendURL   string
	ReverseProxy *httputil.ReverseProxy
	match        func(string) bool
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
func (gateway *Gateway) AddRoute(path, backendURL string, cmp func(string, string) bool, secure bool) {
	backend, err := url.Parse(backendURL)
	if err != nil {
		log.Fatalf("Failed to parse backend URL: %s", err)
	}
	proxy := httputil.NewSingleHostReverseProxy(backend)
	route := &Route{
		Path:         path,
		BackendURL:   backendURL,
		ReverseProxy: proxy,
		match: func(url string) bool {
			return cmp(url, path)
		},
	}
	if secure {
		gateway.SecureRoutes = append(gateway.Routes, route)
	} else {
		gateway.Routes = append(gateway.Routes, route)
	}
}

// ServeHTTP handles incoming HTTP requests
func (gateway *Gateway) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range gateway.Routes {
		if route.match(r.URL.Path) {
			route.ReverseProxy.ServeHTTP(w, r)
			return
		}
	}
	for _, route := range gateway.SecureRoutes {
		if route.match(r.URL.Path) {
			user, err := gateway.authenticate(r)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), userString, user)
			r = r.WithContext(ctx)
			route.ReverseProxy.ServeHTTP(w, r)
			return
		}
	}
	http.NotFound(w, r)
}

// authenticate performs authentication using gRPC metadata
func (gateway *Gateway) authenticate(r *http.Request) (*model.User, error) {
	// Extract user credentials from gRPC metadata
	md, ok := metadata.FromIncomingContext(r.Context())
	if !ok {
		return nil, ErrNoMetadata
	}
	token := md.Get("authorization")[0]

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
		return model.UserFromProto(resp.User), nil
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

	// Routes
	gateway.AddRoute("/user", userService, strings.EqualFold, false)
	gateway.AddRoute("/auth", userService, strings.HasPrefix, false)
	gateway.AddRoute("/notification", notificationService, strings.EqualFold, true)

	http.Handle("/", gateway)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
