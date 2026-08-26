package models

import "time"

type Ticket struct {
	ID          int       `json:"id"`
	Subject     string    `json:"subject" binding:"required"`
	Description string    `json:"description" binding:"required"`
	CustomerID  int       `json:"customer_id" binding:"required"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Status      string    `json:"status"`
	Category    string    `json:"category" binding:"required"`
}
type Representative struct {
	ID        int       `json:"id"`
	Name      string    `json:"name" binding:"max=255"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type RequestTicketMessage struct {
	Message string `json:"message" binding:"required,min=1,max=5000"`
}
type UpdateMessageRequest struct {
	Message string `json:"message" binding:"required,min=1,max=5000"`
}
type TicketMessage struct {
	ID        int       `json:"id"`
	TicketID  int       `json:"ticket_id"`
	UserID    int       `json:"user_id"`
	Message   string    `json:"message" binding:"required,min=1,max=5000"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type Customer struct {
	ID        int       `json:"id"`
	Name      string    `json:"name" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	Email     string    `json:"email" binding:"required,email"`
	UpdatedAt time.Time `json:"updated_at"`
}
type AuthUser struct {
	ID           int       `json:"id"`
	Username     string    `json:"username" binding:"required"`
	PasswordHash string    `json:"passwordHash" binding:"required"`
	Role         string    `json:"role" binding:"required"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	EntityId     *int      `json:"entity_id" binding:"required"`
}
type LoginRequest struct {
	Username string `json:"username" binding:"required,min=5,max=255"`
	Password string `json:"password" binding:"required,min=8,max=255"`
}
type RegisterRequestRepresentative struct {
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required,min=5,max=255"`
	Password string `json:"password" binding:"required,min=8,max=255"`
}
type RegisterRequestCustomer struct {
	Name     string `json:"name" binding:"required,min=1,max=255"`
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required,min=5,max=255"`
	Password string `json:"password" binding:"required,min=8,max=255"`
}
type CreateTicketRequest struct {
	Subject     string `json:"subject"  binding:"required,max=255"`
	Description string `json:"description" binding:"required,min=1,max=5000"`
	Category    string `json:"category" binding:"required"`
}

type ActivityLog struct {
	ID        int       `json:"id"`
	TicketID  int       `json:"ticket_id"`
	UserID    int       `json:"user_id"`
	Action    string    `json:"action"`
	FieldName string    `json:"field_name"`
	OldValue  *string   `json:"old_value"`
	NewValue  *string   `json:"new_value"`
	CreatedAt time.Time `json:"created_at"`
}

type AssignTicketRequest struct {
	TicketID int `json:"ticket_id" binding:"required"`
	RepID    int `json:"rep_id" binding:"required"`
}

const (
	ActionTicketCreated  = "ticket_created"
	ActionTicketUpdated  = "ticket_updated"
	ActionStatusChanged  = "status_changed"
	ActionMessageCreated = "message_created"
	ActionMessageUpdated = "message_updated"
	ActionTicketRevoked  = "ticket_revoked"
	ActionTicketGranted  = "ticket_granted"
	ActionMessageReplied = "message_replied"
)
const (
	CategoryTechnical = "technical"
	CategoryBilling   = "billing"
	CategoryAccount   = "account"
	CategoryBug       = "bug"
	CategoryOther     = "other"
)

type TicketUpdateRequest struct {
	Subject     string `json:"subject" binding:"omitempty,max=255"`
	Description string `json:"description" binding:"omitempty,max=5000"`
	Category    string `json:"category" binding:"omitempty"`
}
type PasswordUpdateRequest struct {
	OldPassword string `json:"old_password" `
	NewPassword string `json:"new_password" `
}
type NotificationRequest struct {
	EventID         string    `json:"event_id"`
	TicketID        int       `json:"ticket_id"`
	ActorUserID     int       `json:"actor_user_id"`
	RecipientUserID int       `json:"recipient_user_id"`
	Type            string    `json:"type"`
	Message         string    `json:"message"`
	OccurredAt      time.Time `json:"occurred_at"`
}
