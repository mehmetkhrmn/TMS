package models

import "time"

type Notification struct {
	ID        int       `json:"id"`
	TicketID  int       `json:"ticket_id" binding:"required"`
	UserID    int       `json:"user_id" binding:"required"`
	Type      string    `json:"type" binding:"required"`
	Message   string    `json:"message" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
}

const (
	NotificationTicketCreated  = "ticket_created"
	NotificationTicketUpdated  = "ticket_updated"
	NotificationMessageCreated = "message_created"
	NotificationMessageUpdated = "message_updated"
	NotificationMessageReplied = "message_replied"
)

func IsValidNotificationType(notificationType string) bool {
	switch notificationType {
	case NotificationTicketCreated,
		NotificationTicketUpdated,
		NotificationMessageCreated,
		NotificationMessageReplied,
		NotificationMessageUpdated:
		return true
	default:
		return false
	}
}
