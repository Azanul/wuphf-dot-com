package postgres

import (
	"context"
	"database/sql"
	"os"
	"strings"
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
	sqlBytes, err := os.ReadFile("../../../pkg/model/user.sql")
	if err != nil {
		return nil, err
	}

	sql := string(sqlBytes)
	statements := strings.Split(sql, ";")
	for _, statement := range statements {
		if strings.TrimSpace(statement) == "" {
			continue
		}
		_, err := db.Exec(statement)
		if err != nil && !strings.Contains(err.Error(), "already exists") {
			return nil, err
		}
	}

	return db, nil
}

func teardownTestDB(db *sql.DB) error {
	// Close the database connection
	if err := db.Close(); err != nil {
		return err
	}

	// Reopen the database connection to perform cleanup operations
	db, err := sql.Open("postgres", "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		return err
	}
	defer db.Close()

	// Execute SQL statements to clear the database
	_, err = db.Exec("DROP SCHEMA public CASCADE; CREATE SCHEMA public;")
	if err != nil {
		return err
	}

	return nil
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
