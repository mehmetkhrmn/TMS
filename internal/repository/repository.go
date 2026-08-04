package repository

import (
	"TMS/internal/models"
	"database/sql"
)

type Repository struct {
	Db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{Db: db}
}

func (r *Repository) CreateTicket(ticket *models.Ticket) error {
	query := `INSERT INTO tickets (subject, description, customer_id, customer_email) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at,status` //returning yazdık routesde kullanıcaz
	err := r.Db.QueryRow(query, ticket.Subject, ticket.Description, ticket.CustomerID,                                                                     //QueryRow always returns a non-nil value
		ticket.CustomerEmail).Scan(&ticket.ID, &ticket.CreatedAt, &ticket.UpdatedAt, &ticket.Status) //scan ile diğer otomatik oluşan verileri aldık
	return err
}

func (r *Repository) GetAllWithStatus(wantedStatus string) ([]models.Ticket, error) {
	var tickets []models.Ticket // veriyi eklemek içi array
	query := `SELECT id,subject,description,customer_id,created_at,updated_at,status,customer_email FROM tickets WHERE status=$1`
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

func (r *Repository) CreateAnswer(answer *models.Answer) error {
	query := "INSERT INTO answers(answer,representative_id,ticket_id,answered_at) values($1,$2,$3,$4) RETURNING id, answered_at,updated_at"
	err := r.Db.QueryRow(query, answer.AnswerText, answer.RepID, answer.TicketID, answer.AnsweredAt).Scan(&answer.ID, &answer.AnsweredAt, &answer.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) UpdateAnswer(id int, answer *models.Answer) error {
	query := "UPDATE answers SET answer=$1 ,updated_at=NOW() WHERE id=$2"
	_, err := r.Db.Exec(query, answer.AnswerText, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetAnswer(ticketID int, answerId int) (*models.Answer, error) {
	var answer models.Answer
	query := `SELECT id,answer,representative_id,ticket_id,answered_at,updated_at FROM answers WHERE id=$1 and ticket_id=$2`
	row, err := r.Db.Query(query, answerId, ticketID)
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
	rows, err := r.Db.Query(query, ticketId)
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
			if err == sql.ErrNoRows {
				return []models.Answer{}, nil
			}
			return nil, err
		}
		answers = append(answers, answer)
	}
	return answers, nil

}
func (r *Repository) CreateCustomerTx(tx *sql.Tx, customer *models.Customer) error {
	query := "INSERT INTO customers(name,email) values($1,$2)RETURNING id,created_at"
	err := tx.QueryRow(query, customer.Name, customer.Email).Scan(&customer.ID, &customer.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}
func (r *Repository) CreateCustomer(customer *models.Customer) error {
	query := "INSERT INTO customers(name,email) values($1,$2)RETURNING id,created_at"
	err := r.Db.QueryRow(query, customer.Name, customer.Email).Scan(&customer.ID, &customer.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}
func (r *Repository) GetCustomer(id int) (*models.Customer, error) {
	var customer models.Customer
	query := "SELECT id,name,email,created_at,updated_at FROM customers WHERE id=$1"
	err := r.Db.QueryRow(query, id).Scan(&customer.ID, &customer.Name, &customer.Email, &customer.CreatedAt, &customer.UpdatedAt) //Changed to QueryRow and used Scan
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &customer, nil
}
func (r *Repository) GetCustomers() ([]models.Customer, error) {
	var customers []models.Customer
	query := "SELECT id,name,email,created_at,updated_at FROM customers"
	rows, err := r.Db.Query(query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var customer models.Customer
		rows.Scan(&customer.ID,
			&customer.Name,
			&customer.Email,
			&customer.CreatedAt,
			&customer.UpdatedAt)
		customers = append(customers, customer)
	}
	return customers, nil
}
func (r *Repository) CreateRepresentativeTx(tx *sql.Tx, rep *models.Representative) error {
	query := "INSERT INTO representatives(name) values($1) RETURNING id, created_at"
	err := tx.QueryRow(query, rep.Name).Scan(&rep.ID, &rep.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}
func (r *Repository) CreateRepresentative(rep *models.Representative) error {
	query := "INSERT INTO representatives(name) values($1) RETURNING id, created_at"
	err := r.Db.QueryRow(query, rep.Name).Scan(&rep.ID, &rep.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}
func (r *Repository) GetAllRepresentatives() ([]models.Representative, error) {
	var representatives []models.Representative
	query := "SELECT id,name,created_at,updated_at FROM representatives"
	rows, err := r.Db.Query(query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var representative models.Representative
		rows.Scan(&representative.ID,
			&representative.Name,
			&representative.CreatedAt,
			&representative.UpdatedAt)
		representatives = append(representatives, representative)
	}
	return representatives, nil
}
func (r *Repository) GetRepresentative(id int) (*models.Representative, error) {
	var representative models.Representative
	query := "SELECT id,name,created_at,updated_at FROM representatives where id=$1"
	err := r.Db.QueryRow(query, id).Scan(&representative.ID, &representative.Name, &representative.CreatedAt, &representative.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &representative, nil
}
func (r *Repository) UpdateRepresentative(id int, representative *models.Representative) error {
	query := "UPDATE representatives SET name=$1 ,updated_at=NOW() WHERE id=$2"
	_, err := r.Db.Exec(query, representative.Name, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) UpdateCustomer(id int, customer *models.Customer) (*models.Customer, error) {
	query := "UPDATE customers SET name=COALESCE(NULLIF($1, ''),name) ,email=COALESCE(NULLIF($2, ''),email),updated_at=NOW() WHERE id=$3"
	_, err := r.Db.Exec(query, customer.Name, customer.Email, id)
	if err != nil {
		return nil, err
	}
	return customer, nil
}
func (r *Repository) GetAuthUserByUsername(username string) (*models.AuthUser, error) {
	var authUser models.AuthUser
	query := "SELECT id,username,password_hash FROM auth_users WHERE username=$1"
	err := r.Db.QueryRow(query, username).Scan(&authUser.ID, &authUser.Username, &authUser.PasswordHash)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &authUser, nil
}
func (r *Repository) CreateAuthUser(tx *sql.Tx, user *models.AuthUser) error {
	query := "INSERT INTO auth_users (username, password_hash,role,entity_id) VALUES ($1, $2,$3,$4) RETURNING id"
	err := tx.QueryRow(query, user.Username, user.PasswordHash, user.Role, user.EntityId).Scan(&user.ID)
	if err != nil {
		return err
	}
	return nil
}
