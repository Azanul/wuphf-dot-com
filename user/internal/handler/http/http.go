package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"reflect"

	"wuphf.com/user/internal/controller/user"
	"wuphf.com/user/internal/repository"
)

// Handler defines a user HTTP handler.
type Handler struct {
	ctrl *user.Controller
}

// New creates a new user HTTP handler.
func New(ctrl *user.Controller) *Handler {
	return &Handler{ctrl}
}

// User handles POST and GET /user requests.
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
		email := req.FormValue("email")
		password := req.FormValue("password")
		if m, err = h.ctrl.Post(req.Context(), email, password); err == nil {
			w.WriteHeader(http.StatusCreated)
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
