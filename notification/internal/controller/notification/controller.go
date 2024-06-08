package notification

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
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
	AssociateUserWithChat(ctx context.Context, userId, chatId string)
	ListChats(ctx context.Context, userId string) ([]string, error)
	ListUsers(ctx context.Context, chatId string) ([]string, error)
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

// Create new chat
func (c *Controller) PostChat(ctx context.Context, sender string, receivers []string) string {
	receivers = append(receivers, sender)
	chatId := generateChatID(receivers)

	c.repo.AssociateUserWithChat(ctx, sender, chatId)
	for _, receiver := range receivers {
		if sender == receiver {
			continue
		}
		c.repo.AssociateUserWithChat(ctx, receiver, chatId)
	}

	return chatId
}

// Post new notification
func (c *Controller) Post(ctx context.Context, sender, chatId, msg string) (string, error) {
	var receivers []string
	var err error
	if chatId == "" {
		chatId = generateChatID([]string{sender, sender})
		receivers = []string{sender}
		c.repo.AssociateUserWithChat(ctx, sender, chatId)
	} else {
		receivers, err = c.repo.ListUsers(ctx, chatId)
		if err != nil {
			return "", err
		}
	}

	for _, receiver := range receivers {
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

		_, err = c.repo.Post(ctx, chatId, notification)
		if err != nil && errors.Is(err, repository.ErrNotFound) {
			return "", repository.ErrNotFound
		}
	}
	return chatId, err
}

// Get returns notification by id
func (c *Controller) Get(ctx context.Context, id string) (*model.Notification, error) {
	res, err := c.repo.Get(ctx, id)
	if err != nil && errors.Is(err, repository.ErrNotFound) {
		return nil, repository.ErrNotFound
	}
	return res, err
}

// List returns list of notifications by chat id
func (c *Controller) List(ctx context.Context, chatId string) ([]*model.Notification, error) {
	res, err := c.repo.List(ctx, chatId)
	if err != nil && errors.Is(err, repository.ErrNotFound) {
		return nil, repository.ErrNotFound
	}
	return res, err
}

// List returns list of chat ids by user id
func (c *Controller) ListChats(ctx context.Context, userId string) ([]string, error) {
	res, err := c.repo.ListChats(ctx, userId)
	if err != nil && errors.Is(err, repository.ErrNotFound) {
		return nil, repository.ErrNotFound
	}
	return res, err
}

// List returns list of user ids by chat id
func (c *Controller) ListUsers(ctx context.Context, chatId string) ([]string, error) {
	res, err := c.repo.ListUsers(ctx, chatId)
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
