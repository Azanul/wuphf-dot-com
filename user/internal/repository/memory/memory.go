package memory

import (
	"context"
	"sync"

	"wuphf.com/user/internal/repository"
	"wuphf.com/user/pkg/model"
)

// Repository defines a memory user repository.
type Repository struct {
	sync.RWMutex
	data     map[string]*model.User
	emailMap map[string]string
}

// New creates a new memory repository
func New() *Repository {
	return &Repository{data: map[string]*model.User{}, emailMap: map[string]string{}}
}

// Post adds a new user
func (r *Repository) Post(_ context.Context, user *model.User) error {
	r.Lock()
	defer r.Unlock()
	r.data[user.ID] = user
	r.emailMap[user.Email] = user.ID
	return nil
}

// Get user by id
func (r *Repository) Get(_ context.Context, id string) (*model.User, error) {
	r.RLock()
	defer r.RUnlock()
	m, ok := r.data[id]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return m, nil
}

// Get user id by email
func (r *Repository) GetIDbyEmail(_ context.Context, email string) (string, error) {
	r.RLock()
	defer r.RUnlock()
	id, ok := r.emailMap[email]
	if !ok {
		return "", repository.ErrNotFound
	}
	return id, nil
}
