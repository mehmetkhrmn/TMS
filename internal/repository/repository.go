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

// ticket ekleme
func (r *Repository) CreateTicket(ticket *models.Ticket) error {
	query := `INSERT INTO tickets (subject, description, customer_id, customer_email, is_done) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at` //returning yazdık routesde kullanıcaz
	err := r.db.QueryRow(query, ticket.Subject, ticket.Description, ticket.CustomerID,
		ticket.CustomerEmail, ticket.IsDone).Scan(&ticket.ID, &ticket.CreatedAt, &ticket.UpdatedAt) //scan ile diğer otomatik oluşan verileri aldık
	return err
}

// TODO: açık olan ticketleri döndür
func (r *Repository) GetAllOpen() ([]models.Ticket, error) {
	var tickets []models.Ticket // veriyi eklemek içi array
	query := `SELECT id,subject,description,customer_id,created_at,updated_at,is_done,customer_email FROM tickets WHERE is_done=FALSE`
	rows, err := r.db.Query(query) //rows u burada cursor olarak kullandık
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
			&ticket.IsDone,
			&ticket.CustomerEmail,
		); err != nil {
			return nil, err
		}
		tickets = append(tickets, ticket) //ticketi arraya ekledik
	}
	return tickets, nil
}

// TODO: Spesifik bir ticketi döndür
func (r *Repository) GetTicket(id int) (*models.Ticket, error) {
	var ticket models.Ticket
	query := `SELECT id,subject,description,customer_id,created_at,updated_at,is_done,customer_email FROM tickets WHERE id=$1`
	row, err := r.db.Query(query, id)
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
		&ticket.IsDone,
		&ticket.CustomerEmail,
	); err != nil {
		return nil, err
	}
	return &ticket, nil
}

func (r *Repository) SetDone(id int, isDone bool) error {
	query := "UPDATE tickets SET is_done=$1,updated_at=NOW() WHERE id=$2"
	_, err := r.db.Exec(query, isDone, id) // döndrüelcek ilk parametre nul olacağı çin boş bıraktık
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) Update(id int, ticket *models.Ticket) error {
	query := "UPDATE tickets SET subject=$1 , description=$2 ,is_done=$3 ,customer_email=$4 ,updated_at=NOW()WHERE id=$5"
	_, err := r.db.Exec(query, ticket.Subject, ticket.Description, ticket.IsDone, ticket.CustomerEmail, id)
	if err != nil {
		return err
	}
	return nil
}

//todo :getalltickets
func (r *Repository)	GetAllTickets() ([]models.Ticket, error) {
	var tickets []models.Ticket
	query:="SELECT id,subject,description,customer_id,created_at,updated_at,is_done,customer_email FROM tickets"
	rows, err := r.db.Query(query)
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
			&ticket.IsDone,
			&ticket.CustomerEmail,
		); err != nil {
			return nil, err
		}
		tickets = append(tickets, ticket) //ticketi arraya ekledik
	}
	return tickets, nil
}