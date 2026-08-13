package repository

import (
	"TMS/notification-service/internal/models"
	"database/sql"
)

func (r *Repository) CreateNotification(notification *models.Notification) error {
	query := "INSERT INTO notifications (ticket_id,user_id,type,message) VALUES($1,$2,$3,$4) RETURNING id,created_at"
	err := r.Db.QueryRow(query, notification.TicketID, notification.UserID, notification.Type, notification.Message).Scan(&notification.ID, &notification.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}
func (r *Repository) GetAllNotifications() ([]models.Notification, error) {
	query := "SELECT id,ticket_id,user_id,type,message,created_at FROM notifications"
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
			&notification.UserID,
			&notification.Type,
			&notification.Message,
			&notification.CreatedAt,
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
	query := "SELECT id,ticket_id,user_id,type,message,created_at FROM notifications WHERE id = $1"
	var notification models.Notification
	err := r.Db.QueryRow(query, notificationId).Scan(&notification.ID, &notification.TicketID, &notification.UserID, &notification.Type, &notification.Message, &notification.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &notification, nil

}
