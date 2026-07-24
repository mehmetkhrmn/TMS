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
	query := `INSERT INTO tickets (subject, description, customer_id, customer_email) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at,status` //returning yazdık routesde kullanıcaz
	err := r.db.QueryRow(query, ticket.Subject, ticket.Description, ticket.CustomerID,
		ticket.CustomerEmail).Scan(&ticket.ID, &ticket.CreatedAt, &ticket.UpdatedAt, &ticket.Status) //scan ile diğer otomatik oluşan verileri aldık
	return err
}

func (r *Repository) GetAllWithStatus(wantedStatus string) ([]models.Ticket, error) {
	var tickets []models.Ticket // veriyi eklemek içi array
	query := `SELECT id,subject,description,customer_id,created_at,updated_at,status,customer_email FROM tickets WHERE status=$1`
	rows, err := r.db.Query(query, wantedStatus) //rows u burada cursor olarak kullandık
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

// TODO: get all close
func (r *Repository) GetAllClosed() ([]models.Ticket, error) {
	var tickets []models.Ticket // veriyi eklemek içi array
	query := `SELECT id,subject,description,customer_id,created_at,updated_at,status,customer_email FROM tickets WHERE status="closed" or status="resolved"`
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
			&ticket.Status,
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
	query := `SELECT id,subject,description,customer_id,created_at,updated_at,status,customer_email FROM tickets WHERE id=$1`
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
		&ticket.Status,
		&ticket.CustomerEmail,
	); err != nil {
		return nil, err
	}
	return &ticket, nil
}

func (r *Repository) SetStatus(id int, status string) error {
	query := "UPDATE tickets SET status=$1,updated_at=NOW() WHERE id=$2"
	_, err := r.db.Exec(query, status, id) // döndrüelcek ilk parametre nul olacağı çin boş bıraktık
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) UpdateTicket(id int, ticket *models.Ticket) error {
	query := `
        UPDATE tickets SET subject = COALESCE(NULLIF($1, ''), subject),description = COALESCE(NULLIF($2, ''), description),status = COALESCE(NULLIF($3, ''), status),customer_email = COALESCE(NULLIF($4, ''), customer_email),updated_at = NOW() WHERE id = $5` // burada eğer değer null ise önceki değeri koru çünkü bize put edilen jsonda her alan dolu olmayabilir
	_, err := r.db.Exec(query, ticket.Subject, ticket.Description, ticket.Status, ticket.CustomerEmail, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetAllTickets() ([]models.Ticket, error) {
	var tickets []models.Ticket
	query := "SELECT id,subject,description,customer_id,created_at,updated_at,status,customer_email FROM tickets"
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
			&ticket.Status,
			&ticket.CustomerEmail,
		); err != nil {
			return nil, err
		}
		tickets = append(tickets, ticket) //ticketi arraya ekledik
	}
	return tickets, nil
}

// todo :create answer check customer id and rep. id
func (r *Repository) CreateAnswer(answer *models.Answer) error {
	query := "INSERT INTO answers(answer,representative_id,ticket_id,answered_at) values($1,$2,$3,$4) RETURNING id, answered_at"
	err := r.db.QueryRow(query, answer.AnswerText, answer.RepID, answer.TicketID, answer.AnsweredAt).Scan(&answer.ID, &answer.AnsweredAt)
	if err != nil {
		return err
	}
	return nil
}

// todo :edit answer
func (r *Repository) UpdateAnswer(id int, answer *models.Answer) error {
	query := "UPDATE answers SET answer=$1 ,updated_at=NOW() WHERE id=$2"
	_, err := r.db.Exec(query, answer.AnswerText, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetAnswer(ticketID int, answerId int) (*models.Answer, error) {
	var answer models.Answer
	query := `SELECT id,answer,representative_id,ticket_id,answered_at,updated_at FROM answers WHERE id=$1 and ticket_id=$2`
	row, err := r.db.Query(query, answerId, ticketID)
	defer row.Close()
	if err != nil {
		return nil, err
	}
	row.Next() //boş tablodan başlıyor
	if err := row.Scan(
		&answer.ID,
		&answer.AnswerText,
		&answer.RepID,
		&answer.TicketID,
		&answer.AnsweredAt,
		&answer.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &answer, nil
}
func (r *Repository) GetAnswers(ticketId int) ([]models.Answer, error) {
	var answers []models.Answer
	query := "select id,answer,representative_id,ticket_id,answered_at,updated_at from answers where ticket_id=$1"
	rows, err := r.db.Query(query, ticketId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var answer models.Answer
		if err := rows.Scan(
			&answer.ID,
			&answer.AnswerText,
			&answer.RepID,
			&answer.TicketID,
			&answer.AnsweredAt,
			&answer.UpdatedAt,
		); err != nil {
			return nil, err
		}
		answers = append(answers, answer)
	}
	return answers, nil

}
