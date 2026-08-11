package repository

import (
	"TMS/internal/models"
	"database/sql"
	"fmt"
	"strings"
)

func (r *Repository) CreateTicket(ticket *models.Ticket) error {
	query := "INSERT INTO tickets (subject, description, customer_id,category) VALUES ($1, $2, $3,$4) RETURNING id, created_at, updated_at,status"                             //returning yazdık routesde kullanıcaz
	err := r.Db.QueryRow(query, ticket.Subject, ticket.Description, ticket.CustomerID, ticket.Category).Scan(&ticket.ID, &ticket.CreatedAt, &ticket.UpdatedAt, &ticket.Status) //QueryRow always returns a non-nil value
	//scan ile diğer otomatik oluşan verileri aldık
	return err
}

func (r *Repository) GetAllWith(wantedStatus string, wantedCategory string) ([]models.Ticket, error) {
	var tickets []models.Ticket

	query := "SELECT id, subject, description, customer_id,created_at, updated_at, status, category FROM tickets"

	var conditions []string
	var args []interface{}
	argID := 1

	if wantedStatus != "" {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argID))
		args = append(args, wantedStatus)
		argID++
	}

	if wantedCategory != "" {
		conditions = append(conditions, fmt.Sprintf("category = $%d", argID))
		args = append(args, wantedCategory)
		argID++
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	rows, err := r.Db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var ticket models.Ticket

		if err := rows.Scan(
			&ticket.ID,
			&ticket.Subject,
			&ticket.Description,
			&ticket.CustomerID,
			&ticket.CreatedAt,
			&ticket.UpdatedAt,
			&ticket.Status,
			&ticket.Category,
		); err != nil {
			return nil, err
		}

		tickets = append(tickets, ticket)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tickets, nil
}

func (r *Repository) GetTicket(id int) (*models.Ticket, error) {
	var ticket models.Ticket
	query := `SELECT id,subject,description,customer_id,created_at,updated_at,status,category FROM tickets WHERE id=$1`
	row, err := r.Db.Query(query, id)
	defer row.Close()
	if err != nil {
		return nil, err
	}
	row.Next() //boş tablodan başlıyor
	if err := row.Scan(
		&ticket.ID,
		&ticket.Subject,
		&ticket.Description,
		&ticket.CustomerID,
		&ticket.CreatedAt,
		&ticket.UpdatedAt,
		&ticket.Status,
		&ticket.Category,
	); err != nil {
		return nil, err
	}
	return &ticket, nil
}

func (r *Repository) SetStatus(id int, status string) error {
	query := "UPDATE tickets SET status=$1,updated_at=NOW() WHERE id=$2"
	_, err := r.Db.Exec(query, status, id) // döndrüelcek ilk parametre nul olacağı çin boş bıraktık Exec executes a query without returning any rows.
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) UpdateTicket(id int, ticket *models.Ticket) error {
	query := `
        UPDATE tickets SET subject = COALESCE(NULLIF($1, ''), subject),description = COALESCE(NULLIF($2, ''), description),status = COALESCE(NULLIF($3, ''), status),category=COALESCE(NULLIF($4,''))updated_at = NOW() WHERE id = $5 RETURNING created_at,updated_at` // burada eğer değer null ise önceki değeri koru çünkü bize put edilen jsonda her alan dolu olmayabilir
	err := r.Db.QueryRow(query, ticket.Subject, ticket.Description, ticket.Status, ticket.Category, id).Scan(&ticket.CreatedAt, &ticket.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetAllTickets() ([]models.Ticket, error) {
	var tickets []models.Ticket
	query := "SELECT id,subject,description,customer_id,created_at,updated_at,status,category FROM tickets"
	rows, err := r.Db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var ticket models.Ticket //fetchlenen veriyi tutmak için
		if err := rows.Scan(
			&ticket.ID,
			&ticket.Subject,
			&ticket.Description,
			&ticket.CustomerID,
			&ticket.CreatedAt,
			&ticket.UpdatedAt,
			&ticket.Status,
			&ticket.Category,
		); err != nil {
			if err == sql.ErrNoRows {
				return []models.Ticket{}, nil
			}
			return nil, err
		}
		tickets = append(tickets, ticket) //ticketi arraya ekledik
	}
	return tickets, nil
}
func (r *Repository) GetCustomerTicket(ticket_id int, cust_id int) (*models.Ticket, error) {
	var ticket models.Ticket
	query := "SELECT ticket_id,subject,description,created_at,updated_at,status,category FROM tickets WHERE tickets.customer_id=$1 and tickets.id=$2"
	err := r.Db.QueryRow(query, cust_id, ticket_id).Scan(
		&ticket.ID,
		&ticket.Subject,
		&ticket.Description,
		&ticket.CreatedAt,
		&ticket.UpdatedAt,
		&ticket.Status,
		&ticket.Category,
	)
	if err != nil {
		return nil, err
	}
	return &ticket, nil
}
func (r *Repository) GetAdminTicket(ticket_id int) (*models.Ticket, error) {
	var ticket models.Ticket
	query := "SELECT id,subject,description,created_at,updated_at,status,customer_id,category FROM tickets WHERE tickets.id=$1"
	row := r.Db.QueryRow(query, ticket_id)
	err := row.Scan(
		&ticket.ID,
		&ticket.Subject,
		&ticket.Description,
		&ticket.CreatedAt,
		&ticket.UpdatedAt,
		&ticket.Status,
		&ticket.CustomerID,
		&ticket.Category,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &ticket, nil
}
func (r *Repository) GetRepresentativeTickets(representative_id int) ([]models.Ticket, error) {
	var tickets []models.Ticket
	query := "SELECT t.id,t.subject,t.description,t.created_at, t.updated_at,t.status,t.customer_id, t.category FROM tickets t JOIN ticket_assignments ta ON ta.ticket_id = t.id WHERE ta.representative_id = $1"

	rows, err := r.Db.Query(query, representative_id)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var ticket models.Ticket
		if err := rows.Scan(
			&ticket.ID,
			&ticket.Subject,
			&ticket.Description,
			&ticket.CreatedAt,
			&ticket.UpdatedAt,
			&ticket.Status,
			&ticket.CustomerID,
			&ticket.Category,
		); err != nil {
			if err == sql.ErrNoRows {
				return []models.Ticket{}, nil
			}
			return nil, err
		}
		tickets = append(tickets, ticket)
	}
	return tickets, nil
}
func (r *Repository) IsTicketOwner(ticketId int, custId int) (bool, error) {
	query := "SELECT EXISTS (SELECT 1  FROM tickets    WHERE id = $1     AND customer_id = $2);"
	var exists bool
	err := r.Db.QueryRow(query, ticketId, custId).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
func (r *Repository) IsUserAssignedTicket(ticketID int, userID int) (bool, error) {
	query := "SELECT EXISTS (SELECT 1 FROM ticket_assignments ta JOIN representatives rep ON rep.id = ta.representative_id JOIN auth_users au ON au.entity_id = rep.id WHERE ta.ticket_id = $1 AND au.id = $2 AND au.role = 'representative')"

	var exists bool

	err := r.Db.QueryRow(query, ticketID, userID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *Repository) GetTicketHistory(ticketId int) ([]models.ActivityLog, error) {
	query := "SELECT id,ticket_id,user_id,action,field_name,old_value,new_value,created_at FROM activity_logs WHERE ticket_id=$1"
	rows, err := r.Db.Query(query, ticketId)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	var logs []models.ActivityLog
	for rows.Next() {
		var log models.ActivityLog
		if err := rows.Scan(&log.ID,
			&log.TicketID,
			&log.UserID,
			&log.Action,
			&log.FieldName,
			&log.OldValue,
			&log.NewValue,
			&log.CreatedAt,
		); err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}
	return logs, nil
}

func (r *Repository) IsRepresentativeAssigned(ticketId int, representativeId int) (bool, error) {
	query := "SELECT EXISTS (    SELECT 1   FROM ticket_assignments   WHERE ticket_id = $1      AND representative_id = $2);"
	var exists bool
	err := r.Db.QueryRow(query, ticketId, representativeId).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
func (r *Repository) AssignRepresentative(ticketId int, representativeId int) error {
	query := "INSERT INTO ticket_assignments(ticket_id,representative_id) VALUES ($1,$2);"
	_, err := r.Db.Exec(query, ticketId, representativeId)
	if err != nil {
		return err
	}
	return nil
}
func (r *Repository) GetCustomerTickets(cust_id int) ([]models.Ticket, error) {
	var tickets []models.Ticket
	query := "SELECT id,subject,description,created_at,updated_at,status,customer_id,category FROM tickets WHERE customer_id=$1; "
	rows, err := r.Db.Query(query, cust_id)
	if err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var ticket models.Ticket
		if err := rows.Scan(
			&ticket.ID,
			&ticket.Subject,
			&ticket.Description,
			&ticket.CreatedAt,
			&ticket.UpdatedAt,
			&ticket.Status,
			&ticket.CustomerID,
			&ticket.Category,
		); err != nil {
			return nil, err
		}
		tickets = append(tickets, ticket)
	}
	return tickets, nil
}
func (r *Repository) GetTicketsByCategory(category string) ([]models.Ticket, error) {
	var tickets []models.Ticket
	query := "SELECT id,subject,description,created_at,updated_at,status,customer_id,category FROM tickets WHERE category=$1"
	rows, err := r.Db.Query(query, category)
	if err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var ticket models.Ticket
		if err := rows.Scan(
			&ticket.ID,
			&ticket.Subject,
			&ticket.Description,
			&ticket.CreatedAt,
			&ticket.UpdatedAt,
			&ticket.Status,
			&ticket.CustomerID,
			&ticket.Category,
		); err != nil {
			return nil, err
		}
		tickets = append(tickets, ticket)

	}
	return tickets, nil
}
