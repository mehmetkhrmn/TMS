package routers

import (
	handlers2 "TMS/ticket-service/internal/handlers"
	"TMS/ticket-service/internal/messaging"
	middleware2 "TMS/ticket-service/internal/middleware"
	"TMS/ticket-service/internal/notification"
	"TMS/ticket-service/internal/repository"
	"os"

	"database/sql"

	"github.com/gin-gonic/gin"
)

func SetupRouter(db *sql.DB, rabbit *messaging.RabbitMQ) *gin.Engine {

	router := gin.Default()
	gin.SetMode(gin.DebugMode)
	repo := repository.NewRepository(db)

	notificationClient := notification.NewClient(
		os.Getenv("NOTIFICATION_SERVICE_URL"),
	)

	authorized := router.Group("/") //route için alt grup oluşturuyorum
	authorized.Use(middleware2.AuthMiddleware())

	admin := authorized.Group("/")
	admin.Use(middleware2.AdminAuthMiddleware())

	customer := authorized.Group("/")
	customer.Use(middleware2.CustomerMiddleware())

	representative := authorized.Group("/")
	representative.Use(middleware2.RepresentativeMiddleware())

	router.POST("/login", func(context *gin.Context) {
		handlers2.Login(context, repo)
	})
	admin.POST("/tickets/:ticket_id/representatives/:representative_id",
		func(c *gin.Context) {
			handlers2.AssignRepresentative(c, repo)
		},
	)
	admin.DELETE("/tickets/:ticket_id/representatives/:representative_id",
		func(c *gin.Context) {
			handlers2.UnassignRepresentative(c, repo)
		})
	authorized.PUT("tickets/:ticket_id/messages/:message_id", func(context *gin.Context) {
		handlers2.UpdateMessage(context, repo, notificationClient, rabbit)
	})

	authorized.POST("/tickets/:ticket_id/messages", func(context *gin.Context) {
		handlers2.CreateMessage(context, repo, notificationClient, rabbit)
	})

	authorized.GET("/tickets/:ticket_id/messages", func(context *gin.Context) {
		handlers2.GetMessages(context, repo)
	})

	authorized.GET("/tickets/:ticket_id/history", func(context *gin.Context) {
		handlers2.GetTicketHistory(context, repo)
	})
	authorized.GET("/tickets/:ticket_id/messages/:message_id", func(context *gin.Context) {
		handlers2.GetMessage(context, repo)
	})

	representative.PUT("/tickets/:ticket_id", func(c *gin.Context) {
		handlers2.UpdateTicket(c, repo, notificationClient, rabbit)
	})

	//statusa göre ticket döndür
	admin.GET("/tickets/", func(context *gin.Context) {
		handlers2.GetTicketWith(context, repo)
	})

	//tek bir ticket görüntüleme
	authorized.GET("/tickets/:ticket_id", func(context *gin.Context) {
		handlers2.GetTicket(context, repo)
	})

	//set status
	representative.PATCH("/tickets/:ticket_id", func(context *gin.Context) {
		handlers2.SetTicketStatus(context, repo, notificationClient, rabbit)
	})

	//bütün ticketleri almak için

	authorized.GET("/tickets", func(context *gin.Context) {
		handlers2.GetAllTickets(context, repo)
	})
	//oluşturmak için
	customer.POST("/tickets", func(context *gin.Context) {
		handlers2.CreateTicket(context, repo, notificationClient, rabbit)
	})

	admin.GET("/customers", func(context *gin.Context) {
		handlers2.GetCustomers(context, repo)
	})

	router.POST("/customers", func(context *gin.Context) {
		handlers2.Register(context, repo)
	})

	representative.GET("/customers/:customer_id", func(context *gin.Context) {
		handlers2.GetCustomer(context, repo)
	})

	representative.PUT("/customers/:customers_id", func(context *gin.Context) {
		handlers2.UpdateCustomer(context, repo)
	})

	admin.GET("/representatives", func(context *gin.Context) {
		handlers2.GetAllRepresentatives(context, repo)
	})

	admin.POST("/representatives", func(context *gin.Context) { //sadece admin temsilci oluşturabilir
		handlers2.Register(context, repo)
	})
	authorized.PUT("/auth/password", func(context *gin.Context) {
		handlers2.UpdatePassword(context, repo)
	})
	admin.PUT("/representatives/:representatives_id", func(context *gin.Context) {
		handlers2.UpdateRepresentative(context, repo)
	})

	admin.GET("/representatives/:representatives_id", func(context *gin.Context) {
		handlers2.GetRepresentative(context, repo)
	})

	return router
}
