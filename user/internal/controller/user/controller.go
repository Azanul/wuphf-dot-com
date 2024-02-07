package user

import (
	"context"
	"errors"

	"wuphf.com/user/internal/repository"
	"wuphf.com/user/pkg/auth"
	"wuphf.com/user/pkg/model"
)

type userRepository interface {
	Get(ctx context.Context, id string) (*model.User, error)
	Post(ctx context.Context, user *model.User) error
	GetIDbyEmail(ctx context.Context, email string) (string, error)
}

// Controller defines a user service controller.
type Controller struct {
	repo userRepository
}

// New creates a user service controller.
func New(repo userRepository) *Controller {
	return &Controller{repo}
}

// Post new user
func (c *Controller) Post(ctx context.Context, email, password string) (string, error) {
	user, err := model.NewUser(email, password)
	if err != nil {
		return "", err
	}
	_, err = c.repo.GetIDbyEmail(ctx, email)
	if err == nil {
		return "", repository.ErrDuplicate
	}
	err = c.repo.Post(ctx, user)
	if err != nil && errors.Is(err, repository.ErrNotFound) {
		return "", repository.ErrNotFound
	}
	return user.ID, err
}

// Get returns user by id.
func (c *Controller) Get(ctx context.Context, id string) (*model.User, error) {
	res, err := c.repo.Get(ctx, id)
	if err != nil && errors.Is(err, repository.ErrNotFound) {
		return nil, repository.ErrNotFound
	}
	return res, err
}

// Login new user
func (c *Controller) Login(ctx context.Context, email, password string) (string, error) {
	id, err := c.repo.GetIDbyEmail(ctx, email)
	if err != nil {
		return "", repository.ErrNotFound
	}

	user, err := c.repo.Get(ctx, id)
	if err != nil {
		return "", err
	}
	hashed_password, err := model.HashPassword(password)
	if err != nil {
		return "", err
	}
	if user.Password == hashed_password {
		return "", repository.ErrInvalidCredentials
	}

	token, err := auth.GenerateToken(id)
	if err != nil {
		return "", err
	}

	return token, nil
}
