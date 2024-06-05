package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/Azanul/wuphf-dot-com/user/gen"
	"github.com/Azanul/wuphf-dot-com/user/internal/controller/user"
	grpchandler "github.com/Azanul/wuphf-dot-com/user/internal/handler/grpc"
	httphandler "github.com/Azanul/wuphf-dot-com/user/internal/handler/http"

	"google.golang.org/grpc"

	"github.com/Azanul/wuphf-dot-com/user/internal/repository/memory"
)

func main() {
	log.Println("Starting the user service")
	repo := memory.New()
	ctrl := user.New(repo)

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
