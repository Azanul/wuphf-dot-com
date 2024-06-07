package memory

import (
	"context"
	"strconv"
	"sync"

	"github.com/Azanul/wuphf-dot-com/notification/internal/repository"
	"github.com/Azanul/wuphf-dot-com/notification/pkg/model"
)

// Repository defines a memory notification repository
type Repository struct {
	sync.RWMutex
	data      map[string][]*model.Notification
	userChats map[string][]string
	chatUsers map[string][]string
}

// New creates a new memory repository
func New() *Repository {
	return &Repository{
		data:      map[string][]*model.Notification{},
		userChats: map[string][]string{},
		chatUsers: map[string][]string{},
	}
}

// Post adds a new notification
func (r *Repository) Post(_ context.Context, chatID string, n *model.Notification) (int, error) {
	r.Lock()
	defer r.Unlock()
	r.data[chatID] = append(r.data[chatID], n)
	return len(r.data[chatID]) - 1, nil
}

// Get notification by id
func (r *Repository) Get(_ context.Context, id string) (*model.Notification, error) {
	r.RLock()
	defer r.RUnlock()
	n, ok := r.data[id[:repository.ID_LENGTH]]
	if !ok {
		return nil, repository.ErrNotFound
	}
	idx, err := strconv.Atoi(string(id[repository.ID_LENGTH:]))
	if err != nil {
		return nil, err
	}
	if idx >= len(n) {
		return nil, repository.ErrNotFound
	}
	return n[idx], nil
}

// List notification by chat id
func (r *Repository) List(_ context.Context, chatID string) ([]*model.Notification, error) {
	r.RLock()
	defer r.RUnlock()
	n_list, ok := r.data[chatID]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return n_list, nil
}

// AssociateUserWithChat associates a user with a chat
func (r *Repository) AssociateUserWithChat(_ context.Context, userID, chatID string) {
	r.Lock()
	defer r.Unlock()
	r.userChats[userID] = append(r.userChats[userID], chatID)
	r.chatUsers[chatID] = append(r.chatUsers[chatID], userID)
}

// ListChats retrieves chat ids for a given user id
func (r *Repository) ListChats(_ context.Context, userID string) ([]string, error) {
	r.RLock()
	defer r.RUnlock()
	chatIDs, ok := r.userChats[userID]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return chatIDs, nil
}

// ListUsers retrieves user ids for a given chat id
func (r *Repository) ListUsers(_ context.Context, chatID string) ([]string, error) {
	r.RLock()
	defer r.RUnlock()
	userIDs, ok := r.chatUsers[chatID]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return userIDs, nil
}
