package postgres

import (
	"context"
	"database/sql"
	"log"
	"testing"

	_ "github.com/lib/pq"

	"github.com/Azanul/wuphf-dot-com/notification/pkg/model"
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

func TestRepository(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Error setting up test database: %v\n", err)
	}
	defer teardownTestDB(db)

	repo := New(db)

	ctx := context.Background()

	// Test Post
	t.Run("TestPost", func(t *testing.T) {
		chatID := "test_chat_id"
		n := &model.Notification{Msg: "test message"}
		id, err := repo.Post(ctx, chatID, n)
		if err != nil {
			t.Errorf("Error posting notification: %v\n", err)
		}

		// Assuming the first inserted id is 1
		if id != 1 {
			t.Errorf("Expected id to be 1, got %d\n", id)
		}
	})

	// Test Get
	t.Run("TestGet", func(t *testing.T) {
		expectedMessage := "test message"
		id := "1" // Assuming the id of the previously inserted notification is 1
		n, err := repo.Get(ctx, id)
		if err != nil {
			t.Errorf("Error getting notification: %v\n", err)
		}

		if n.Msg != expectedMessage {
			t.Errorf("Expected message to be %s, got %s\n", expectedMessage, n.Msg)
		}
	})

	// Test List
	t.Run("TestList", func(t *testing.T) {
		chatID := "test_chat_id"
		expectedMessages := []string{"test message"}
		notifications, err := repo.List(ctx, chatID)
		if err != nil {
			t.Errorf("Error listing notifications: %v\n", err)
		}

		if len(notifications) != len(expectedMessages) {
			t.Errorf("Expected %d notifications, got %d\n", len(expectedMessages), len(notifications))
		}

		for i, n := range notifications {
			if n.Msg != expectedMessages[i] {
				t.Errorf("Expected message to be %s, got %s\n", expectedMessages[i], n.Msg)
			}
		}
	})
}
