package repository

import (
	"TMS/internal/models"
	"database/sql"
)

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
func (r *Repository) UpdateCustomer(id int, customer *models.Customer) (*models.Customer, error) {
	query := "UPDATE customers SET name=COALESCE(NULLIF($1, ''),name) ,email=COALESCE(NULLIF($2, ''),email),updated_at=NOW() WHERE id=$3"
	_, err := r.Db.Exec(query, customer.Name, customer.Email, id)
	if err != nil {
		return nil, err
	}
	return customer, nil
}
