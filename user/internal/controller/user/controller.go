package user

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"net/textproto"
	"reflect"

	"github.com/Azanul/wuphf-dot-com/user/internal/repository"
	"github.com/Azanul/wuphf-dot-com/user/pkg/auth"
	"github.com/Azanul/wuphf-dot-com/user/pkg/model"
	"github.com/IBM/sarama"
)

type userRepository interface {
	Get(ctx context.Context, id string) (*model.User, error)
	Post(ctx context.Context, user *model.User) error
	GetIDbyEmail(ctx context.Context, email string) (string, error)
}

// Controller defines a user service controller
type Controller struct {
	repo          userRepository
	kafkaProducer sarama.AsyncProducer
}

// New creates a user service controller
func New(repo userRepository, kafkaProducer sarama.AsyncProducer) *Controller {
	return &Controller{repo, kafkaProducer}
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

	message, err := createMultipartFormData(map[string]any{"sender": user.ID, "receivers": []string{user.ID}})
	if err != nil {
		log.Fatalf("Failed to create multipart form data: %v", err)
	}
	c.kafkaProducer.BeginTxn()
	c.kafkaProducer.Input() <- &sarama.ProducerMessage{
		Topic: "chats",
		Value: sarama.StringEncoder(message),
	}

	message, err = createMultipartFormData(map[string]any{"sender": user.ID, "chat_id": "", "msg": "Wuphf"})
	if err != nil {
		log.Fatalf("Failed to create multipart form data: %v", err)
	}
	c.kafkaProducer.Input() <- &sarama.ProducerMessage{
		Topic: "notifications",
		Value: sarama.StringEncoder(message),
	}

	return user.ID, err
}

// Get returns user by id
func (c *Controller) Get(ctx context.Context, id string) (*model.User, error) {
	res, err := c.repo.Get(ctx, id)
	if err != nil && errors.Is(err, repository.ErrNotFound) {
		return nil, repository.ErrNotFound
	}
	return res, err
}

// Login new user
func (c *Controller) Login(ctx context.Context, email, password string) (string, string, error) {
	id, err := c.repo.GetIDbyEmail(ctx, email)
	if err != nil {
		return "", "", repository.ErrNotFound
	}

	user, err := c.repo.Get(ctx, id)
	if err != nil {
		return "", "", err
	}
	hashed_password, err := model.HashPassword(password)
	if err != nil {
		return "", "", err
	}
	if user.Password == hashed_password {
		return "", "", repository.ErrInvalidCredentials
	}

	token, err := auth.GenerateToken(id)
	if err != nil {
		return "", "", err
	}

	return user.ID, token, nil
}

func createMultipartFormData(fields map[string]any) (string, error) {
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)
	if err := writer.SetBoundary("--------------------------"); err != nil {
		return "", err
	}

	for key, val := range fields {
		h := make(textproto.MIMEHeader)
		h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"`, key))

		v := reflect.ValueOf(val)
		switch v.Kind() {
		case reflect.String:
			part, err := writer.CreatePart(h)
			if err != nil {
				return "", err
			}
			part.Write([]byte(val.(string)))
		case reflect.Slice:
			for i := 0; i < v.Len(); i++ {
				part, err := writer.CreatePart(h)
				if err != nil {
					return "", err
				}
				part.Write([]byte(v.Index(i).Interface().(string)))
			}
		default:
			return "", fmt.Errorf("unsupported type for field %s", key)
		}
	}

	err := writer.Close()
	if err != nil {
		return "", err
	}

	return b.String(), nil
}
