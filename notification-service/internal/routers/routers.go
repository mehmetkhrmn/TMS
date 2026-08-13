package routers

import (
	"TMS/notification-service/internal/handlers"
	"TMS/notification-service/internal/repository"
	"database/sql"

	"github.com/gin-gonic/gin"
)

func SetupRouter(db *sql.DB) *gin.Engine {
	router := gin.Default()
	notifications := router.Group("/notifications")
	repo := repository.NewRepository(db)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "notification service running",
		})
	})

	notifications.POST("", func(context *gin.Context) {
		handlers.CreateNotification(context, repo)
	})
	notifications.GET("", func(context *gin.Context) {
		handlers.GetNotifications(context, repo)
	})
	notifications.GET("/:id", func(context *gin.Context) {
		handlers.GetNotification(context, repo)
	})
	return router

}
