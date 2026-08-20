package models

import "time"

type Notification struct {
	EventID         string    `json:"event_id" binding:"required"`
	ID              int       `json:"id"`
	TicketID        int       `json:"ticket_id" binding:"required"`
	RecipientUserID int       `json:"recipient_user_id"`
	ActorUserID     *int       `json:"actor_user_id"`// sistemin kendisi de olabilir
	Type            string    `json:"type" binding:"required"`
	Message         string    `json:"message" binding:"required"`
	OccurredAt      time.Time `json:"occurred_at" binding:"required"`
	CreatedAt       time.Time `json:"created_at"`
}

const (
	NotificationTicketCreated  = "ticket_created"
	NotificationTicketUpdated  = "ticket_updated"
	NotificationMessageCreated = "message_created"
	NotificationMessageUpdated = "message_updated"
	NotificationMessageReplied = "message_replied"
	NotificationTicketAssigned = "ticket_assigned"
	NotificationTicketRevoked  = "ticket_revoked"
)

func IsValidNotificationType(notificationType string) bool {
	switch notificationType {
	case NotificationTicketCreated,
		NotificationTicketUpdated,
		NotificationMessageCreated,
		NotificationMessageReplied,
		NotificationTicketAssigned,
		NotificationTicketRevoked,
		NotificationMessageUpdated:
		return true
	default:
		return false
	}
}
