package memory

import (
	"context"
	"testing"

	"github.com/Azanul/wuphf-dot-com/notification/internal/repository"

	"github.com/Azanul/wuphf-dot-com/notification/pkg/model"
)

func TestMemoryRepository(t *testing.T) {
	repo := New()

	ctx := context.Background()
	expectedNotification, _ := model.NewNotification("sender1", "receiver1", "testBody", "")
	chatID := repository.RandStringBytesMaskImpr(repository.ID_LENGTH)
	userID := "user1"

	// Test adding a notification
	t.Run("TestPostNotification", func(t *testing.T) {
		_, err := repo.Post(ctx, chatID, expectedNotification)
		if err != nil {
			t.Errorf("Error posting notification: %v", err)
		}
	})

	// Test getting a notification
	t.Run("TestGetNotification", func(t *testing.T) {
		ctx := context.Background()
		notification, err := repo.Get(ctx, chatID+"0")
		if err != nil {
			t.Errorf("Error getting notification: %v", err)
		}
		if notification.Sender != expectedNotification.Sender ||
			notification.Receiver != expectedNotification.Receiver ||
			notification.Msg != expectedNotification.Msg {
			t.Errorf("Expected notification %v, got %v", expectedNotification, notification)
		}
	})

	// Test listing notifications
	t.Run("TestListNotifications", func(t *testing.T) {
		ctx := context.Background()
		expectedNotifications := []*model.Notification{expectedNotification}
		notifications, err := repo.List(ctx, string(chatID))
		if err != nil {
			t.Errorf("Error listing notifications: %v", err)
		}
		if len(notifications) != len(expectedNotifications) {
			t.Errorf("Expected %d notifications, got %d", len(expectedNotifications), len(notifications))
		}
		for i, notification := range notifications {
			if notification.Sender != expectedNotifications[i].Sender ||
				notification.Receiver != expectedNotifications[i].Receiver ||
				notification.Msg != expectedNotifications[i].Msg {
				t.Errorf("Expected notification %v, got %v", expectedNotifications[i], notification)
			}
		}
	})

	// Test associating a user with a chat
	t.Run("TestAssociateUserWithChat", func(t *testing.T) {
		repo.AssociateUserWithChat(ctx, userID, chatID)
		chatIDs, err := repo.ListChats(context.Background(), userID)
		if err != nil {
			t.Errorf("Error getting chat IDs: %v", err)
		}
		if len(chatIDs) != 1 || chatIDs[0] != chatID {
			t.Errorf("Expected chat ID %s for user %s, got %v", chatID, userID, chatIDs)
		}
	})

	// Test retrieving chat IDs by user ID
	t.Run("TestGetChatIDsByUserID", func(t *testing.T) {
		anotherChatID := repository.RandStringBytesMaskImpr(repository.ID_LENGTH)
		repo.AssociateUserWithChat(ctx, userID, anotherChatID)

		chatIDs, err := repo.ListChats(context.Background(), userID)
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
		repo.AssociateUserWithChat(ctx, anotherUserID, chatID)

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
