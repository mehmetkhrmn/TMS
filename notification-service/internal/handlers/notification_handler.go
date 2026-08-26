package handlers

import (
	"TMS/notification-service/internal/repository"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetNotification(c *gin.Context, repo *repository.Repository) {
	notifID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	role := c.GetString("role")
	userId := int(c.GetFloat64("user_id"))
	if role != "admin" {
		ok, err := repo.IsRecipientOfNotification(userId, notifID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "authorization error"})
			return
		}
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization is required"})
			return
		}
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
	role := c.GetString("role")
	userId := int(c.GetFloat64("user_id"))
	if role != "admin" {
		notifs, err := repo.GetNotificationByUserId(userId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		c.JSON(http.StatusOK, notifs)
		return
	}
	notifs, err := repo.GetAllNotifications()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, notifs)
}
