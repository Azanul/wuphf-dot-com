package postgres

import (
	"context"
	"database/sql"

	"github.com/Azanul/wuphf-dot-com/user/internal/repository"
	"github.com/Azanul/wuphf-dot-com/user/pkg/model"
)

// UserRepository defines a PostgreSQL user repository
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new PostgreSQL user repository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Post adds a new user
func (r *UserRepository) Post(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO users (id, email, password) VALUES ($1, $2, $3)
	`
	_, err := r.db.ExecContext(ctx, query, user.ID, user.Email, user.Password)
	return err
}

// Get retrieves a user by id
func (r *UserRepository) Get(ctx context.Context, id string) (*model.User, error) {
	query := `
		SELECT id, email, password FROM users WHERE id = $1
	`
	row := r.db.QueryRowContext(ctx, query, id)
	user := &model.User{}
	err := row.Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return user, nil
}

// GetIDByEmail retrieves a user id by email
func (r *UserRepository) GetIDByEmail(ctx context.Context, email string) (string, error) {
	query := `
		SELECT id FROM users WHERE email = $1
	`
	var id string
	err := r.db.QueryRowContext(ctx, query, email).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", repository.ErrNotFound
		}
		return "", err
	}
	return id, nil
}
