package models

import "time"

type Ticket struct {
	ID            int       `json:"id"`
	Subject       string    `json:"subject"`
	Description   string    `json:"description"`
	CustomerID    int       `json:"customer_id"`
	CustomerEmail string    `json:"customer_email"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Status        string    `json:"status"`
}

type Representative struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Answer struct {
	ID         int       `json:"id"`
	AnswerText string    `json:"answer"`
	RepID      int       `json:"representative_id"`
	TicketID   int       `json:"ticket_id"`
	AnsweredAt time.Time `json:"answered_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
type Customer struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	Email     string    `json:"email"`
	UpdatedAt time.Time `json:"updated_at"`
}
type AuthUser struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"passwordHash"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	EntityId     int       `json:"entity_id"`
}
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type RegisterRequestRepresentative struct {
	Name     string `json:"name"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" `
	EntityID int    `json:"entity_id" `
}
type RegisterRequestCustomer struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" `
	EntityID int    `json:"entity_id"`
}
