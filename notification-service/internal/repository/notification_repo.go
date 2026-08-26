package repository

import (
	"TMS/notification-service/internal/models"
	"database/sql"
	"errors"
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
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
	defer func() {
		_ = rows.Close()
	}()
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
			if errors.Is(err, sql.ErrNoRows) {
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
	err := r.Db.QueryRow(query, notificationId).Scan(&notification.ID, &notification.TicketID, &notification.ActorUserID, &notification.RecipientUserID, &notification.Type, &notification.Message, &notification.CreatedAt, &notification.OccurredAt, &notification.EventID)
	if err != nil {
		return nil, err
	}
	return &notification, nil

}
func (r *Repository) GetNotificationByUserId(userId int) ([]models.Notification, error) {
	query := "SELECT id,ticket_id,actor_user_id,recipient_user_id,type,message,created_at,occurred_at,event_id FROM notifications WHERE notifications.recipient_user_id = $1"
	rows, err := r.Db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
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
			if errors.Is(err, sql.ErrNoRows) {
				return []models.Notification{}, nil
			}
			return nil, err
		}
		notifications = append(notifications, notification)
	}
	return notifications, nil
}
func (r *Repository) IsRecipientOfNotification(userId int, notifId int) (bool, error) {
	query := "SELECT EXISTS (SELECT 1  FROM notifications    WHERE recipient_user_id = $1     AND id = $2);"
	var exists bool
	err := r.Db.QueryRow(query, userId, notifId).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
