package postgres

import (
	"context"
	"database/sql"
	"log"
	"testing"

	_ "github.com/lib/pq"

	"github.com/Azanul/wuphf-dot-com/user/pkg/model"
)

func setupTestDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func teardownTestDB(db *sql.DB) {
	if err := db.Close(); err != nil {
		log.Fatalf("Error closing test database: %v\n", err)
	}
}

func TestUserRepository(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Error setting up test database: %v\n", err)
	}
	defer teardownTestDB(db)

	repo := NewUserRepository(db)

	ctx := context.Background()

	hashed_password, err := model.HashPassword("password")
	if err != nil {
		t.Errorf("Error hashing password: %v\n", err)
	}
	user := &model.User{
		ID:       "test_user_id",
		Email:    "test@example.com",
		Password: hashed_password,
	}

	// Test Post & Get
	t.Run("TestPost", func(t *testing.T) {
		err = repo.Post(ctx, user)
		if err != nil {
			t.Errorf("Error posting user: %v\n", err)
		}

		// Now try to retrieve the same user
		retrievedUser, err := repo.Get(ctx, user.ID)
		if err != nil {
			t.Errorf("Error retrieving user: %v\n", err)
		}

		if retrievedUser.ID != user.ID || retrievedUser.Email != user.Email || retrievedUser.Password != user.Password {
			t.Errorf("Retrieved user does not match original user: %v != %v", user, retrievedUser)
		}
	})

	// Test GetIDByEmail
	t.Run("TestGetIDByEmail", func(t *testing.T) {
		id, err := repo.GetIDByEmail(ctx, user.Email)
		if err != nil {
			t.Errorf("Error retrieving user id by email: %v\n", err)
		}

		if id != user.ID {
			t.Errorf("Retrieved user ID does not match original user ID")
		}
	})
}
