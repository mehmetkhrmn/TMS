package repository

import (
	"TMS/internal/models"
	"database/sql"
)

func (r *Repository) CreateTicket(ticket *models.Ticket) error {
	query := `INSERT INTO tickets (subject, description, customer_id, customer_email) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at,status` //returning yazdık routesde kullanıcaz
	err := r.Db.QueryRow(query, ticket.Subject, ticket.Description, ticket.CustomerID,                                                                     //QueryRow always returns a non-nil value
		ticket.CustomerEmail).Scan(&ticket.ID, &ticket.CreatedAt, &ticket.UpdatedAt, &ticket.Status) //scan ile diğer otomatik oluşan verileri aldık
	return err
}

func (r *Repository) GetAllWithStatus(wantedStatus string) ([]models.Ticket, error) {
	var tickets []models.Ticket // veriyi eklemek içi array
	query := `SELECT id,subject,description,customer_id,created_at,updated_at,status,customer_email FROM tickets WHERE status=$1`
	if wantedStatus == "" {
		query = `SELECT id,subject,description,customer_id,created_at,updated_at,status,customer_email FROM tickets`
	}
	rows, err := r.Db.Query(query, wantedStatus) //Query executes a query that returns rows
	if err != nil {
		//hata kontrolü
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
			&ticket.CustomerEmail,
		); err != nil {
			return nil, err
		}
		tickets = append(tickets, ticket) //ticketi arraya ekledik
	}
	return tickets, nil
}

func (r *Repository) GetTicket(id int) (*models.Ticket, error) {
	var ticket models.Ticket
	query := `SELECT id,subject,description,customer_id,created_at,updated_at,status,customer_email FROM tickets WHERE id=$1`
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
		&ticket.CustomerEmail,
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
        UPDATE tickets SET subject = COALESCE(NULLIF($1, ''), subject),description = COALESCE(NULLIF($2, ''), description),status = COALESCE(NULLIF($3, ''), status),customer_email = COALESCE(NULLIF($4, ''), customer_email),updated_at = NOW() WHERE id = $5 RETURNING created_at,updated_at` // burada eğer değer null ise önceki değeri koru çünkü bize put edilen jsonda her alan dolu olmayabilir
	err := r.Db.QueryRow(query, ticket.Subject, ticket.Description, ticket.Status, ticket.CustomerEmail, id).Scan(&ticket.CreatedAt, &ticket.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetAllTickets() ([]models.Ticket, error) {
	var tickets []models.Ticket
	query := "SELECT id,subject,description,customer_id,created_at,updated_at,status,customer_email FROM tickets"
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
			&ticket.CustomerEmail,
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
	query := "SELECT ticket_id,subject,description,created_at,updated_at,status FROM tickets WHERE tickets.customer_id=$1 and tickets.id=$2"
	err := r.Db.QueryRow(query, cust_id, ticket_id).Scan(
		&ticket.ID,
		&ticket.Subject,
		&ticket.Description,
		&ticket.CreatedAt,
		&ticket.UpdatedAt,
		&ticket.Status,
	)
	if err != nil {
		return nil, err
	}
	return &ticket, nil
}
func (r *Repository) GetAdminTicket(ticket_id int) (*models.Ticket, error) {
	var ticket models.Ticket
	query := "SELECT ticket_id,subject,description,created_at,updated_at,status,customer_id,customer_email FROM tickets WHERE tickets.id=$1"
	row := r.Db.QueryRow(query, ticket_id)
	err := row.Scan(
		&ticket.ID,
		&ticket.Subject,
		&ticket.Description,
		&ticket.CreatedAt,
		&ticket.UpdatedAt,
		&ticket.Status,
		&ticket.CustomerID,
		&ticket.CustomerEmail)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &ticket, nil
}
func (r *Repository) GetRepresentativeTickets(customer_id int, representative_id int) ([]models.Ticket, error) {
	var tickets []models.Ticket
	query := "SELECT ticket_id,subject,description,created_at,updated_at,status,customer_id,customer_email FROM tickets WHERE customer_id=$1 and representative_id=$2 "
	rows, err := r.Db.Query(query, customer_id, representative_id)
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
			&ticket.CustomerEmail,
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
