package notification

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/Azanul/wuphf-dot-com/notification/internal/repository"
	"github.com/Azanul/wuphf-dot-com/notification/pkg/model"
)

type notificationIntegration interface {
	Name() string
	Notify(receiver, message string) (string, error)
}

type notificationRepository interface {
	Get(ctx context.Context, id string) (*model.Notification, error)
	Post(ctx context.Context, chatId string, n *model.Notification) (int, error)
	List(ctx context.Context, chatId string) ([]*model.Notification, error)
}

// Controller defines a notification service controller
type Controller struct {
	repo         notificationRepository
	integrations []notificationIntegration
}

// New creates a notification service controller
func New(repo notificationRepository) *Controller {
	return &Controller{repo, []notificationIntegration{}}
}

func (c *Controller) AddIntegration(ni notificationIntegration) {
	c.integrations = append(c.integrations, ni)
}

// Post new notification
func (c *Controller) Post(ctx context.Context, sender, receiver, msg string) (string, error) {
	chatId := generateChatID([]string{sender, receiver})

	reference := map[string]string{}
	for _, i := range c.integrations {
		res, err := i.Notify(receiver, msg)
		if err == nil {
			reference[i.Name()] = res
		} else {
			reference[i.Name()] = err.Error()
		}
	}
	refBytes, err := json.Marshal(reference)
	if err != nil {
		return "", err
	}

	notification, err := model.NewNotification(sender, receiver, msg, string(refBytes))
	if err != nil {
		return "", err
	}

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

// List returns list of notifications by ids of participants
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
