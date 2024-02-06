package main

import (
	"log"
	"net/http"

	"wuphf.com/user/internal/controller/user"
	httphandler "wuphf.com/user/internal/handler/http"
	"wuphf.com/user/internal/repository/memory"
)

func main() {
	log.Println("Starting the user service")
	repo := memory.New()
	ctrl := user.New(repo)
	h := httphandler.New(ctrl)

	// Endpoints
	http.Handle("/user", http.HandlerFunc(h.User))
	http.Handle("/register", http.HandlerFunc(h.Register))
	http.Handle("/login", http.HandlerFunc(h.Login))

	if err := http.ListenAndServe(":8081", nil); err != nil {
		panic(err)
	}
}
