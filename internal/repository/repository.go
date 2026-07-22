package repository

import (
	"TMS/internal/models"
	"database/sql"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}
func (r *Repository) CreateTicket(ticket *models.Ticket) error {
	query := `INSERT INTO tickets (subject, description, customer_id, customer_email, is_done) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at` //returning yazdık routesde kullanıcaz
	err := r.db.QueryRow(query, ticket.Subject, ticket.Description, ticket.CustomerID,
		ticket.CustomerEmail, ticket.IsDone).Scan(&ticket.ID, &ticket.CreatedAt, &ticket.UpdatedAt) //scan ile diğer otomatik oluşan verileri aldık
	return err
}
