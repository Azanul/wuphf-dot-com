package notification

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"wuphf.com/notification/internal/repository"
	"wuphf.com/notification/pkg/model"
)

type notificationRepository interface {
	Get(ctx context.Context, id string) (*model.Notification, error)
	Post(ctx context.Context, chatId string, n *model.Notification) (int, error)
	List(ctx context.Context, chatId string) ([]*model.Notification, error)
}

// Controller defines a notification service controller.
type Controller struct {
	repo notificationRepository
}

// New creates a notification service controller.
func New(repo notificationRepository) *Controller {
	return &Controller{repo}
}

// Post new notification
func (c *Controller) Post(ctx context.Context, sender, receiver, msg string) (string, error) {
	notification, err := model.NewNotification(sender, receiver, msg)
	if err != nil {
		return "", err
	}
	chatId := generateChatID([]string{sender, receiver})
	i, err := c.repo.Post(ctx, chatId, notification)
	if err != nil && errors.Is(err, repository.ErrNotFound) {
		return "", repository.ErrNotFound
	}
	return chatId + strconv.Itoa(i), err
}

// Get returns notification by id
func (c *Controller) Get(ctx context.Context, id string) (*model.Notification, error) {
	res, err := c.repo.Get(ctx, id)
	if err != nil && errors.Is(err, repository.ErrNotFound) {
		return nil, repository.ErrNotFound
	}
	return res, err
}

// List returns list of notifications by ids of participants.
func (c *Controller) List(ctx context.Context, chatId string) ([]*model.Notification, error) {
	res, err := c.repo.List(ctx, chatId)
	if err != nil && errors.Is(err, repository.ErrNotFound) {
		return nil, repository.ErrNotFound
	}
	return res, err
}

// Helper function to generate chat id from the user ids
func generateChatID(userIDs []string) string {
	sort.Strings(userIDs)

	userIDsStr := strings.Join(userIDs, "_")

	hash := sha256.New()
	hash.Write([]byte(userIDsStr))
	chatID := fmt.Sprintf("%x", hash.Sum(nil))

	return chatID
}
