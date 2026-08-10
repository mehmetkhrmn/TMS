package repository

import (
	"TMS/internal/models"
	"database/sql"
)

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

func (r *Repository) GetAnswer(answerId int) (*models.Answer, error) {
	var answer models.Answer
	query := `SELECT id,answer,representative_id,ticket_id,answered_at,updated_at FROM answers WHERE id=$1`
	row, err := r.Db.Query(query, answerId)
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
func (r *Repository) IsAnswerMatchWithRep(answerID, representativeID int) (bool, error) {
	query := "SELECT EXISTS (SELECT 1  FROM answers    WHERE id = $1     AND representative_id = $2);"
	var exists bool
	err := r.Db.QueryRow(query, answerID, representativeID).Scan(&exists)

	if err != nil {
		return false, err
	}
	return exists, nil
}
func (r *Repository) IsAnswerBelongsToTicket(answerID int, ticketID int) (bool, error) {
	query := "SELECT EXISTS (SELECT 1  FROM answers    WHERE id = $1     AND ticket_id = $2)"
	var exists bool
	err := r.Db.QueryRow(query, answerID, ticketID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
func (r *Repository) GetTicketIDByAnswerID(answerID int) (int, error) {
	ticketID := 0
	query := "SELECT id FROM tickets WHERE answer_id=$1"
	err := r.Db.QueryRow(query, answerID).Scan(&ticketID)
	if err != nil {
		return 0, err
	}
	return ticketID, nil
}
