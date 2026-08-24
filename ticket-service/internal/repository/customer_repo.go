package repository

import (
	"TMS/ticket-service/internal/models"
	"database/sql"
	"log"
)

func (r *Repository) CreateCustomerTx(tx *sql.Tx, customer *models.Customer) error {
	query := "INSERT INTO customers(name,email) values($1,$2)RETURNING id,created_at"
	err := tx.QueryRow(query, customer.Name, customer.Email).Scan(&customer.ID, &customer.CreatedAt)
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
		return nil, err
	}

	return &customer, nil
}
func (r *Repository) GetCustomers() ([]models.Customer, error) {
	var customers []models.Customer
	query := "SELECT id,name,email,created_at,updated_at FROM customers"
	rows, err := r.Db.Query(query)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Println(err)
		}
	}()
	for rows.Next() {
		var customer models.Customer
		err := rows.Scan(&customer.ID,
			&customer.Name,
			&customer.Email,
			&customer.CreatedAt,
			&customer.UpdatedAt)

		if err != nil {
			return nil, err
		}
		customers = append(customers, customer)
	}
	return customers, nil
}
func (r *Repository) UpdateCustomer(id int, customer *models.Customer) (*models.Customer, error) {
	query := "UPDATE customers SET name=COALESCE(NULLIF($1, ''),name) ,email=COALESCE(NULLIF($2, ''),email),updated_at=NOW() WHERE id=$3"
	_, err := r.Db.Exec(query, customer.Name, customer.Email, id)
	if err != nil {
		return nil, err
	}
	return customer, nil
}
func (r *Repository) GetCustomerEmail(customerID int) (string, error) {
	query := "SELECT email FROM customers WHERE id=$1"
	row := r.Db.QueryRow(query, customerID)
	var email string
	err := row.Scan(&email)
	return email, err
}
func (r *Repository) GetUserIDByCustID(custId int) (int, error) {
	var userID int
	query := "SELECT id FROM auth_users WHERE entity_id=$1 AND role='customer'"
	err := r.Db.QueryRow(query, custId).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}
func (r *Repository) IsCustomerAssignedToRepresentative(customerID int, representativeID int) (bool, error) {
	query := "SELECT EXISTS (SELECT 1 FROM ticket_assignments a JOIN tickets t ON t.id = a.ticket_id WHERE a.representative_id = $1 AND t.customer_id = $2);"
	var assigned bool
	err := r.Db.QueryRow(query, representativeID, customerID).Scan(&assigned)
	if err != nil {
		return false, err
	}
	return assigned, nil
}
