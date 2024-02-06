package memory

import (
	"context"
	"strconv"
	"sync"

	"wuphf.com/notification/internal/repository"
	"wuphf.com/notification/pkg/model"
)

const ID_LENGTH = 64

// Repository defines a memory notification repository.
type Repository struct {
	sync.RWMutex
	data map[string][]*model.Notification
}

// New creates a new memory repository
func New() *Repository {
	return &Repository{data: map[string][]*model.Notification{}}
}

// Post adds a new notification
func (r *Repository) Post(_ context.Context, chatId string, n *model.Notification) (int, error) {
	r.Lock()
	defer r.Unlock()
	r.data[chatId] = append(r.data[chatId], n)
	return len(r.data[chatId]) - 1, nil
}

// Get notification by id
func (r *Repository) Get(_ context.Context, id string) (*model.Notification, error) {
	r.RLock()
	defer r.RUnlock()
	n, ok := r.data[id[:ID_LENGTH]]
	if !ok {
		return nil, repository.ErrNotFound
	}
	idx, err := strconv.Atoi(string(id[ID_LENGTH:]))
	if err != nil {
		return nil, err
	}
	if idx >= len(n) {
		return nil, repository.ErrNotFound
	}
	return n[idx], nil
}

// List notification by user ids list
func (r *Repository) List(_ context.Context, chatId string) ([]*model.Notification, error) {
	r.RLock()
	defer r.RUnlock()
	n_list, ok := r.data[chatId]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return n_list, nil
}
