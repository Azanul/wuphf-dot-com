package main

import (
	"log"
	"net/http"

	"wuphf.com/notification/internal/controller/notification"
	httphandler "wuphf.com/notification/internal/handler/http"
	"wuphf.com/notification/internal/repository/memory"
)

func main() {
	log.Println("Starting the notification service")
	repo := memory.New()
	ctrl := notification.New(repo)
	h := httphandler.New(ctrl)

	// Endpoints
	http.Handle("/notification", http.HandlerFunc(h.Notification))
	http.Handle("/history", http.HandlerFunc(h.History))

	if err := http.ListenAndServe(":8082", nil); err != nil {
		panic(err)
	}
}
