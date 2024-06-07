package postgres

import (
	"context"
	"database/sql"
	"log"
	"os"
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

	sqlBytes, err := os.ReadFile("../../../pkg/model/notification.sql")
	if err != nil {
		return nil, err
	}

	sql := string(sqlBytes)
	_, err = db.Exec(sql)
	if err != nil {
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
	chatID := "test_chat_id"
	userID := "user1"

	// Test Post
	t.Run("TestPost", func(t *testing.T) {
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

	// Test associating a user with a chat
	t.Run("TestAssociateUserWithChat", func(t *testing.T) {
		ctx := context.Background()
		err := repo.AssociateUserWithChat(ctx, userID, chatID)
		if err != nil {
			t.Errorf("Error associating user with chat: %v", err)
		}
		chatIDs, err := repo.ListChats(ctx, userID)
		if err != nil {
			t.Errorf("Error getting chat IDs: %v", err)
		}
		if len(chatIDs) != 1 || chatIDs[0] != chatID {
			t.Errorf("Expected chat ID %s for user %s, got %v", chatID, userID, chatIDs)
		}
	})

	// Test retrieving chat IDs by user ID
	t.Run("TestListChats", func(t *testing.T) {
		ctx := context.Background()
		anotherChatID := "test_chat_id_2"
		err := repo.AssociateUserWithChat(ctx, userID, anotherChatID)
		if err != nil {
			t.Errorf("Error associating user with another chat: %v", err)
		}

		chatIDs, err := repo.ListChats(ctx, userID)
		if err != nil {
			t.Errorf("Error getting chat IDs: %v", err)
		}
		expectedChatIDs := []string{chatID, anotherChatID}
		if len(chatIDs) != len(expectedChatIDs) {
			t.Errorf("Expected %d chat IDs, got %d", len(expectedChatIDs), len(chatIDs))
		}
		for i, id := range chatIDs {
			if id != expectedChatIDs[i] {
				t.Errorf("Expected chat ID %s, got %s", expectedChatIDs[i], id)
			}
		}
	})

	// Test retrieving user IDs by chat ID
	t.Run("TestListUsers", func(t *testing.T) {
		ctx := context.Background()
		anotherUserID := "user2"
		err := repo.AssociateUserWithChat(ctx, anotherUserID, chatID)
		if err != nil {
			t.Errorf("Error associating another user with chat: %v", err)
		}

		userIDs, err := repo.ListUsers(ctx, chatID)
		if err != nil {
			t.Errorf("Error getting user IDs: %v", err)
		}
		expectedUserIDs := []string{userID, anotherUserID}
		if len(userIDs) != len(expectedUserIDs) {
			t.Errorf("Expected %d user IDs, got %d", len(expectedUserIDs), len(userIDs))
		}
		for i, id := range userIDs {
			if id != expectedUserIDs[i] {
				t.Errorf("Expected user ID %s, got %s", expectedUserIDs[i], id)
			}
		}
	})
}
