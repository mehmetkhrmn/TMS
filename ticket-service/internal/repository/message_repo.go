package repository

import (
	"TMS/ticket-service/internal/models"
	"database/sql"
	"log"
)

func (r *Repository) CreateMessage(tx *sql.Tx, message *models.TicketMessage) error {
	query := "INSERT INTO ticket_messages(message,user_id,ticket_id) values($1,$2,$3) RETURNING id, created_at,updated_at"
	err := r.Db.QueryRow(query, message.Message, message.UserID, message.TicketID).Scan(&message.ID, &message.CreatedAt, &message.UpdatedAt)
	if err != nil {
		return err
	}
	err = r.CreateActivityLog(tx, message.TicketID, message.UserID, models.ActionMessageCreated, "message", "", "created")
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) UpdateMessage(tx *sql.Tx, id int, message *models.UpdateMessageRequest) error {
	query := "UPDATE ticket_messages SET message=$1 ,updated_at=NOW() WHERE id=$2"
	_, err := tx.Exec(query, message.Message, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetMessage(messageId int) (*models.TicketMessage, error) {
	var message models.TicketMessage
	query := "SELECT id,message,user_id,ticket_id,created_at,updated_at FROM ticket_messages WHERE id=$1"
	err := r.Db.QueryRow(query, messageId).Scan(&message.ID,
		&message.Message,
		&message.UserID,
		&message.TicketID,
		&message.CreatedAt,
		&message.UpdatedAt)


	if err != nil {
		return nil, err
	}

	return &message, nil
}
func (r *Repository) GetMessages(ticketId int) ([]models.TicketMessage, error) {
	var messages []models.TicketMessage
	query := "select id,message,user_id,ticket_id,created_at,updated_at from ticket_messages where ticket_id=$1"
	rows, err := r.Db.Query(query, ticketId)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Printf("error closing rows: %s", err)
		}
	}()

	for rows.Next() {
		var message models.TicketMessage
		if err := rows.Scan(
			&message.ID,
			&message.Message,
			&message.UserID,
			&message.TicketID,
			&message.CreatedAt,
			&message.UpdatedAt,
		); err != nil {

			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, nil

}
func (r *Repository) IsMessageMatchWithUser(messageID, userID int) (bool, error) {
	query := "SELECT EXISTS (SELECT 1  FROM ticket_messages    WHERE id = $1     AND user_id = $2);"
	var exists bool
	err := r.Db.QueryRow(query, messageID, userID).Scan(&exists)

	if err != nil {
		return false, err
	}
	return exists, nil
}
func (r *Repository) IsMessageBelongsToTicket(messageID int, ticketID int) (bool, error) {
	query := "SELECT EXISTS (SELECT 1  FROM ticket_messages    WHERE id = $1     AND ticket_id = $2)"
	var exists bool
	err :=
		r.Db.QueryRow(query, messageID, ticketID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
func (r *Repository) GetTicketIdByMessageId(messageID int) (int, error) {
	ticketID := 0
	query := "SELECT ticket_id FROM ticket_messages WHERE id=$1"
	err := r.Db.QueryRow(query, messageID).Scan(&ticketID)
	if err != nil {
		return 0, err
	}
	return ticketID, nil
}
