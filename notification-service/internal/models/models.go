package models

import "time"

type Notification struct {
	ID        int       `json:"id"`
	TicketID  int       `json:"ticket_id"`
	UserID    int       `json:"user_id"`
	Type      string    `json:"type"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

const (
	NotificationTicketCreated  = "ticket_created"
	NotificationTicketUpdated  = "ticket_updated"
	NotificationMessageCreated = "message_created"
	NotificationMessageReplied = "message_replied"

)
