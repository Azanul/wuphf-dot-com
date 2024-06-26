package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"reflect"

	"github.com/Azanul/wuphf-dot-com/user/internal/controller/user"
	"github.com/Azanul/wuphf-dot-com/user/internal/repository"
)

// Handler defines a user HTTP handler
type Handler struct {
	ctrl *user.Controller
}

// New creates a new user HTTP handler
func New(ctrl *user.Controller) *Handler {
	return &Handler{ctrl}
}

// User handles POST and GET /user requests
func (h *Handler) User(w http.ResponseWriter, req *http.Request) {
	var err error
	var m any
	ctx := req.Context()

	switch req.Method {
	case http.MethodGet:
		id := req.FormValue("id")
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if m, err = h.ctrl.Get(ctx, id); err == nil {
			w.WriteHeader(http.StatusOK)
		}
	case http.MethodPost:
		var token string
		email := req.FormValue("email")
		password := req.FormValue("password")
		m, token, err = h.ctrl.Post(ctx, email, password)
		if err == nil {
			w.Header().Add("AUTHORIZATION", token)
			w.WriteHeader(http.StatusCreated)
			m = map[string]string{"user_id": m.(string)}
		}
	}

	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
		} else if errors.Is(err, repository.ErrDuplicate) {
			m = "user already exists"
			w.WriteHeader(http.StatusBadRequest)
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

// Login handles POST /register requests
func (h *Handler) Register(w http.ResponseWriter, req *http.Request) {
	var err error
	var m any
	ctx := req.Context()

	switch req.Method {
	case http.MethodPost:
		var token string
		email := req.FormValue("email")
		password := req.FormValue("password")
		m, token, err = h.ctrl.Post(ctx, email, password)
		if err == nil {
			w.Header().Add("AUTHORIZATION", token)
			w.WriteHeader(http.StatusCreated)
			m = map[string]string{"user_id": m.(string)}
		}
	}

	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
		} else if errors.Is(err, repository.ErrDuplicate) {
			m = "user already exists"
			w.WriteHeader(http.StatusBadRequest)
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

// Login handles POST /login requests
func (h *Handler) Login(w http.ResponseWriter, req *http.Request) {
	var err error
	var m any
	ctx := req.Context()

	switch req.Method {
	case http.MethodPost:
		var token string
		email := req.FormValue("email")
		password := req.FormValue("password")
		m, token, err = h.ctrl.Login(ctx, email, password)
		if err == nil {
			w.Header().Add("AUTHORIZATION", token)
			w.WriteHeader(http.StatusOK)
			m = map[string]string{"user_id": m.(string)}
		}
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
