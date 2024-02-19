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
