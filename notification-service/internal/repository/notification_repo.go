package repository

import (
	"TMS/notification-service/internal/models"
	"database/sql"
)

func (r *Repository) CreateNotification(notification *models.Notification) error {
	query := "INSERT INTO notifications (event_id,ticket_id,actor_user_id,recipient_user_id,type,message,occurred_at) VALUES ($1,$2,COALESCE($3,0),COALESCE($4,0),$5,$6,COALESCE($7,NOW())) ON CONFLICT (event_id) DO NOTHING RETURNING id,created_at"
	err := r.Db.QueryRow(
		query,
		notification.EventID,
		notification.TicketID,
		notification.ActorUserID,
		notification.RecipientUserID,
		notification.Type,
		notification.Message,
		notification.OccurredAt,
	).Scan(
		&notification.ID,
		&notification.CreatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}
func (r *Repository) GetAllNotifications() ([]models.Notification, error) {
	query := "SELECT id,ticket_id,actor_user_id,recipient_user_id,type,message,created_at,occurred_at,event_id FROM notifications"
	rows, err := r.Db.Query(query)
	if err != nil {
		return nil, err
	}
	var notifications []models.Notification
	for rows.Next() {
		var notification models.Notification
		if err := rows.Scan(
			&notification.ID,
			&notification.TicketID,
			&notification.ActorUserID,
			&notification.RecipientUserID,
			&notification.Type,
			&notification.Message,
			&notification.CreatedAt,
			&notification.OccurredAt,
			&notification.EventID,
		); err != nil {
			if err == sql.ErrNoRows {
				return []models.Notification{}, nil
			}
			return nil, err
		}
		notifications = append(notifications, notification)
	}
	return notifications, nil
}
func (r *Repository) GetNotification(notificationId int) (*models.Notification, error) {
	query := "SELECT id,ticket_id,actor_user_id,recipient_user_id,type,message,created_at,occurred_at,event_id FROM notifications WHERE id = $1"
	var notification models.Notification
	err := r.Db.QueryRow(query, notificationId).Scan(&notification.ID, &notification.TicketID, &notification.ActorUserID, notification.RecipientUserID, &notification.Type, &notification.Message, &notification.CreatedAt, &notification.OccurredAt, &notification.EventID)
	if err != nil {
		return nil, err
	}
	return &notification, nil

}
