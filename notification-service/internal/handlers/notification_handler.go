package handlers

import (
	"TMS/notification-service/internal/models"
	"TMS/notification-service/internal/repository"
	"database/sql"
	"log/slog"
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

	// geçerli notification typeları buraya yazılmalı


	if !models.IsValidNotificationType(notification.Type) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid notification type"})
		return
	}

	err = repo.CreateNotification(&notification)
	if err != nil {
		slog.Error("notification service unavailable", "error", err)

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
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "notification not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	if notif == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "notification not found"})
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
