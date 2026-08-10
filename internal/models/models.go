package models

import "time"

type Ticket struct {
	ID            int       `json:"id"`
	Subject       string    `json:"subject" binding:"required"`
	Description   string    `json:"description" binding:"required"`
	CustomerID    int       `json:"customer_id" binding:"required"`
	CustomerEmail string    `json:"customer_email" binding:"required,email"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Status        string    `json:"status"`
}
type Representative struct {
	ID        int       `json:"id"`
	Name      string    `json:"name" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Answer struct {
	ID         int       `json:"id"`
	AnswerText string    `json:"answer" binding:"required"`
	RepID      int       `json:"representative_id" binding:"required"`
	TicketID   int       `json:"ticket_id" binding:"required"`
	AnsweredAt time.Time `json:"answered_at"`
	UpdatedAt  time.Time `json:"updated_at"`
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
	EntityId     int       `json:"entity_id" binding:"required"`
}
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type RegisterRequestRepresentative struct {
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type RegisterRequestCustomer struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type CreateTicketRequest struct {
	Subject     string `json:"subject" binding:"required"`
	Description string `json:"description" binding:"required"`
}
type CreateAnswerRequest struct {
	AnswerText string `json:"answer" binding:"required"`
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
	ActionTicketCreated = "ticket_created"
	ActionTicketUpdated = "ticket_updated"
	ActionStatusChanged = "status_changed"
	ActionAnswerCreated = "answer_created"
	ActionAnswerUpdated = "answer_updated"
)
