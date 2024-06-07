package postgres

import (
	"context"
	"database/sql"

	"github.com/Azanul/wuphf-dot-com/notification/internal/repository"
	"github.com/Azanul/wuphf-dot-com/notification/pkg/model"
)

// Repository defines a sql notification repository
type Repository struct {
	db *sql.DB
}

// New creates a new memory repository
func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// Post adds a new notification
func (r *Repository) Post(ctx context.Context, chatID string, n *model.Notification) (int, error) {
	query := `
		INSERT INTO notifications (chat_id, msg)
		VALUES ($1, $2)
		RETURNING id
	`

	var id int
	err := r.db.QueryRowContext(ctx, query, chatID, n.Msg).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// Get notification by id
func (r *Repository) Get(ctx context.Context, id string) (*model.Notification, error) {
	query := `
		SELECT msg FROM notifications WHERE id = $1
	`

	var message string
	err := r.db.QueryRowContext(ctx, query, id).Scan(&message)
	if err == sql.ErrNoRows {
		return nil, repository.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return &model.Notification{Msg: message}, nil
}

// List notification by user chatID
func (r *Repository) List(ctx context.Context, chatID string) ([]*model.Notification, error) {
	query := `
		SELECT msg FROM notifications WHERE chat_id = $1
	`

	rows, err := r.db.QueryContext(ctx, query, chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []*model.Notification
	for rows.Next() {
		var message string
		if err := rows.Scan(&message); err != nil {
			return nil, err
		}
		notifications = append(notifications, &model.Notification{Msg: message})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return notifications, nil
}

// AssociateUserWithChat associates a user with a chat
func (r *Repository) AssociateUserWithChat(ctx context.Context, userID, chatID string) error {
	query := `
		INSERT INTO user_chats (user_id, chat_id)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING
	`
	_, err := r.db.ExecContext(ctx, query, userID, chatID)
	return err
}

// ListChats retrieves chat IDs for a given user ID
func (r *Repository) ListChats(ctx context.Context, userID string) ([]string, error) {
	query := `
		SELECT chat_id FROM user_chats WHERE user_id = $1
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chatIDs []string
	for rows.Next() {
		var chatID string
		if err := rows.Scan(&chatID); err != nil {
			return nil, err
		}
		chatIDs = append(chatIDs, chatID)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return chatIDs, nil
}

// ListUsers retrieves user IDs for a given chat ID
func (r *Repository) ListUsers(ctx context.Context, chatID string) ([]string, error) {
	query := `
		SELECT user_id FROM user_chats WHERE chat_id = $1
	`

	rows, err := r.db.QueryContext(ctx, query, chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userIDs []string
	for rows.Next() {
		var userID string
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}
		userIDs = append(userIDs, userID)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return userIDs, nil
}
