package handlers

import (
	"TMS/ticket-service/internal/messaging"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HealthLive(c *gin.Context) { //service acik
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func HealthReady(c *gin.Context, db *sql.DB, rabbit *messaging.RabbitMQ) { //db baglantisi var
	if err := db.Ping(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "db not ready",
		})
		return
	}
	if rabbit.Conn == nil || rabbit.Conn.IsClosed() {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "mq not ready",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
	})
}
