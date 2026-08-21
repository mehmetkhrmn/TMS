package repository

import (
	"TMS/ticket-service/internal/models"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
)

func (r *Repository) CreateTicket(tx *sql.Tx, ticket *models.Ticket) error {
	query := "INSERT INTO tickets (subject, description, customer_id,category) VALUES ($1, $2, $3,$4) RETURNING id, created_at, updated_at,status"                           //returning yazdık routesde kullanıcaz
	err := tx.QueryRow(query, ticket.Subject, ticket.Description, ticket.CustomerID, ticket.Category).Scan(&ticket.ID, &ticket.CreatedAt, &ticket.UpdatedAt, &ticket.Status) //QueryRow always returns a non-nil value
	if err != nil {
		return err
	}
	//scan ile diğer otomatik oluşan verileri aldık

	return nil

}

func (r *Repository) GetAllWith(
	wantedStatus string,
	wantedCategory string,
	limit int,
	offset int,
) ([]models.Ticket, error) {

	var tickets []models.Ticket

	query := "SELECT id, subject, description, customer_id, created_at, updated_at, status, category FROM tickets"

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

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argID, argID+1)

	args = append(args, limit, offset)

	rows, err := r.Db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Printf("Error closing rows: %v", err)
		}
	}()

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
	err := r.Db.QueryRow(query, id).Scan(
		&ticket.ID,
		&ticket.Subject,
		&ticket.Description,
		&ticket.CustomerID,
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

func (r *Repository) SetStatus(tx *sql.Tx, id int, status string) error {
	query := "UPDATE tickets SET status=$1,updated_at=NOW() WHERE id=$2"
	_, err := tx.Exec(query, status, id) // döndrüelcek ilk parametre nul olacağı çin boş bıraktık Exec executes a query without returning any rows.
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) UpdateTicket(tx *sql.Tx, id int, ticket *models.Ticket) error {
	query := `UPDATE tickets SET subject = COALESCE(NULLIF($1, ''), subject),description = COALESCE(NULLIF($2, ''), description),status = COALESCE(NULLIF($3, ''), status),category=COALESCE(NULLIF($4,''),category) ,updated_at = NOW() WHERE id = $5 RETURNING created_at,updated_at` // burada eğer değer null ise önceki değeri koru çünkü bize put edilen jsonda her alan dolu olmayabilir
	err := tx.QueryRow(query, ticket.Subject, ticket.Description, ticket.Status, ticket.Category, id).Scan(&ticket.CreatedAt, &ticket.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetAllTickets(limit int, offset int) ([]models.Ticket, error) {
	var tickets []models.Ticket

	query := "SELECT id,subject,description,customer_id,created_at,updated_at,status,category FROM tickets ORDER BY created_at DESC LIMIT $1 OFFSET $2"

	rows, err := r.Db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Printf("Error closing rows: %v", err)
		}
	}()

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
func (r *Repository) GetCustomerTicket(ticketId int, custId int) (*models.Ticket, error) {
	var ticket models.Ticket
	query := "SELECT id,subject,description,created_at,updated_at,status,category FROM tickets WHERE tickets.customer_id=$1 and tickets.id=$2"
	err := r.Db.QueryRow(query, custId, ticketId).Scan(
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
func (r *Repository) GetAdminTicket(ticketId int) (*models.Ticket, error) {
	var ticket models.Ticket
	query := "SELECT id,subject,description,created_at,updated_at,status,customer_id,category FROM tickets WHERE tickets.id=$1"
	row := r.Db.QueryRow(query, ticketId)
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
		return nil, err
	}
	return &ticket, nil
}
func (r *Repository) GetRepresentativeTickets(
	representativeID int,
	limit int,
	offset int,
) ([]models.Ticket, error) {

	var tickets []models.Ticket

	query := "SELECT t.id,t.subject,t.description,t.created_at,t.updated_at,t.status,t.customer_id,t.category FROM tickets t JOIN ticket_assignments ta ON ta.ticket_id=t.id WHERE ta.representative_id=$1 ORDER BY t.created_at DESC LIMIT $2 OFFSET $3"

	rows, err := r.Db.Query(query, representativeID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Printf("Error closing rows: %v", err)
		}
	}()

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

	if err := rows.Err(); err != nil {
		return nil, err
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

	if err != nil {
		return nil, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Printf("Error closing rows: %v", err)
		}
	}()
	var logs []models.ActivityLog
	for rows.Next() {
		var aLog models.ActivityLog
		if err := rows.Scan(&aLog.ID,
			&aLog.TicketID,
			&aLog.UserID,
			&aLog.Action,
			&aLog.FieldName,
			&aLog.OldValue,
			&aLog.NewValue,
			&aLog.CreatedAt,
		); err != nil {
			return nil, err
		}
		logs = append(logs, aLog)
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
func (r *Repository) AssignRepresentative(tx *sql.Tx, ticketId int, representativeId int) error {
	query := "INSERT INTO ticket_assignments(ticket_id,representative_id) VALUES ($1,$2);"
	_, err := tx.Exec(query, ticketId, representativeId)
	if err != nil {
		return err
	}

	return nil
}
func (r *Repository) GetCustomerTickets(
	custId int,
	limit int,
	offset int,
) ([]models.Ticket, error) {

	var tickets []models.Ticket

	query := "SELECT id, subject, description, created_at, updated_at, status, customer_id, category FROM tickets WHERE customer_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3"

	rows, err := r.Db.Query(query, custId, limit, offset)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Printf("Error closing rows: %v", err)
		}
	}()

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

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tickets, nil
}

func (r *Repository) UnAssignRepresentative(tx *sql.Tx, ticketId int, repId int) error {
	query := "DELETE FROM ticket_assignments WHERE ticket_id=$1 AND representative_id=$2"
	res, err := tx.Exec(query, ticketId, repId)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()

	if rows == 0 {
		return errors.New("ticket_id is invalid")
	}
	if err != nil {
		return err
	}

	return nil

}
func (r *Repository) GetCustomerUserIDByTicketID(ticketID int) (int, error) {
	var userID int

	query := "SELECT au.id FROM tickets t JOIN auth_users au ON au.entity_id = t.customer_id WHERE t.id = $1 AND au.role = 'customer'"

	err := r.Db.QueryRow(query, ticketID).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (r *Repository) GetRepresentativeUserIDByTicketID(ticketID int) (int, error) {
	var userID int

	query := "SELECT au.id FROM ticket_assignments ta JOIN auth_users au ON au.entity_id = ta.representative_id WHERE ta.ticket_id = $1 AND au.role = 'representative'"

	err := r.Db.QueryRow(query, ticketID).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}
