package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"reflect"

	"github.com/Azanul/wuphf-dot-com/notification/internal/controller/notification"
	"github.com/Azanul/wuphf-dot-com/notification/internal/repository"
)

// Handler defines a notification HTTP handler
type Handler struct {
	ctrl *notification.Controller
}

// New creates a new notification HTTP handler
func New(ctrl *notification.Controller) *Handler {
	return &Handler{ctrl}
}

// Notify handles POST and GET /notification requests
func (h *Handler) Notification(w http.ResponseWriter, req *http.Request) {
	var err error
	var m any

	ctx := req.Context()

	switch req.Method {
	case http.MethodGet:
		id := req.FormValue("id")
		if m, err = h.ctrl.Get(ctx, id); err == nil {
			w.WriteHeader(http.StatusOK)
		}
	case http.MethodPost:
		sender := req.FormValue("sender")
		receiver := req.FormValue("receiver")
		msg := req.FormValue("msg")
		if m, err = h.ctrl.Post(ctx, sender, receiver, msg); err == nil {
			w.WriteHeader(http.StatusCreated)
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			log.Printf("Repository get error: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
	if m != nil && !(reflect.ValueOf(m).Kind() == reflect.Ptr && reflect.ValueOf(m).IsNil()) && m != "" {
		if err := json.NewEncoder(w).Encode(m); err != nil {
			log.Printf("Response encode error: %v\n", err)
		}
	}
}

// History handles GET /history requests
func (h *Handler) History(w http.ResponseWriter, req *http.Request) {
	var err error
	var m any

	switch req.Method {
	case http.MethodGet:
		id := req.FormValue("chatId")
		if id == "" {
			id = req.FormValue("userId")
			if m, err = h.ctrl.ListChats(req.Context(), id); err == nil {
				w.WriteHeader(http.StatusOK)
			}
			m = []map[string]any{{"chatId": m, "messages": []string{}}}
		} else {
			if m, err = h.ctrl.List(req.Context(), id); err == nil {
				w.WriteHeader(http.StatusOK)
			}
		}

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			log.Printf("Repository get error: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
	if m != nil && !(reflect.ValueOf(m).Kind() == reflect.Ptr && reflect.ValueOf(m).IsNil()) && m != "" {
		if err := json.NewEncoder(w).Encode(m); err != nil {
			log.Printf("Response encode error: %v\n", err)
		}
	}
}
