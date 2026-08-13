package handlers

import (
	"TMS/notification-service/internal/models"
	"TMS/notification-service/internal/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateNotification(c *gin.Context, repo *repository.Repository) {
	var notification models.Notification
	err := c.ShouldBindJSON(&notification)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	switch notification.Type {	// geçerli notification typeları buraya yazılmalı
	case models.NotificationTicketCreated,
		models.NotificationTicketUpdated,
		models.NotificationMessageCreated,
		models.NotificationMessageReplied:

	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid notification type",
		})
		return
	}
	err = repo.CreateNotification(&notification)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, notification)

}
func GetNotification(c *gin.Context, repo *repository.Repository) {
	notifID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	notif, err := repo.GetNotification(notifID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, notif)
}
func GetNotifications(c *gin.Context, repo *repository.Repository) {
	notifs, err := repo.GetAllNotifications()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, notifs)
}
