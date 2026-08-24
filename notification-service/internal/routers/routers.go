package routers

import (
	"TMS/notification-service/internal/handlers"
	"TMS/notification-service/internal/messaging"
	"TMS/notification-service/internal/repository"

	"github.com/gin-gonic/gin"
)

func SetupRouter(repo *repository.Repository, rabbit *messaging.RabbitMQ) *gin.Engine {
	router := gin.Default()
	notifications := router.Group("/notifications")
	router.GET("/health/live", func(c *gin.Context) {
		handlers.HealthLive(c)
	})
	router.GET("/health/ready", func(c *gin.Context) {
		handlers.HealthReady(c, repo, rabbit)
	})

	notifications.GET("", func(context *gin.Context) {
		handlers.GetNotifications(context, repo)
	})
	notifications.GET("/:id", func(context *gin.Context) {
		handlers.GetNotification(context, repo)
	})
	return router

}
