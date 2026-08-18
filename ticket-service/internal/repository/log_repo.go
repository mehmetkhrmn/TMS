package repository

import "database/sql"

func (r *Repository) CreateActivityLog(
	tx *sql.Tx,
	ticketID int,
	userID int,
	action string,
	fieldName string,
	oldValue string,
	newValue string,
) error {

	query := "INSERT INTO activity_logs(ticket_id, user_id, action, field_name, old_value, new_value)VALUES ($1, $2, $3, $4, $5, $6)"

	_, err := tx.Exec(
		query,
		ticketID,
		userID,
		action,
		fieldName,
		oldValue,
		newValue,
	)

	return err
}
