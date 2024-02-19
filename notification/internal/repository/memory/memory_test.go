package memory

import (
	"context"
	"testing"

	"github.com/Azanul/wuphf-dot-com/notification/internal/repository"

	"github.com/Azanul/wuphf-dot-com/notification/pkg/model"
)

func TestMemoryRepository(t *testing.T) {
	repo := New()

	expectedNotification, _ := model.NewNotification("sender1", "receiver1", "testBody")
	chatID := repository.RandStringBytesMaskImpr(ID_LENGTH)

	// Test adding a notification
	t.Run("TestPostNotification", func(t *testing.T) {
		ctx := context.Background()
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
}
