package routers

import (
	"TMS/internal/handlers"
	"TMS/internal/middleware"
	"TMS/internal/repository"

	"database/sql"

	"github.com/gin-gonic/gin"
)

func SetupRouter(db *sql.DB) *gin.Engine {

	router := gin.Default()
	gin.SetMode(gin.DebugMode)
	repo := repository.NewRepository(db)

	authorized := router.Group("/") //route için alt grup oluşturuyorum
	authorized.Use(middleware.AuthMiddleware())

	admin := authorized.Group("/")
	admin.Use(middleware.AdminAuthMiddleware())

	representative := authorized.Group("/")
	representative.Use(middleware.RepresentativeMiddleware())

	router.POST("/login", func(context *gin.Context) {
		handlers.Login(context, repo)
	})
	admin.POST("/tickets/:ticket_id/representatives/:representative_id",
		func(c *gin.Context) {
			handlers.AssignRepresentative(c, repo)
		},
	)
	//:TODO buraya eksta ticket id ekledik ama formaliteden gibi eksta chechk koy
	authorized.PUT("tickets/:ticket_id/messages/:message_id", func(context *gin.Context) {
		handlers.UpdateMessage(context, repo)
	})

	authorized.POST("/tickets/:ticket_id/messages", func(context *gin.Context) {
		handlers.CreateMessage(context, repo)
	})

	authorized.GET("/tickets/:ticket_id/messages", func(context *gin.Context) {
		handlers.GetMessages(context, repo)
	})

	authorized.GET("/tickets/:ticket_id/history", func(context *gin.Context) {
		handlers.GetTicketHistory(context, repo)
	})
	authorized.GET("/tickets/:ticket_id/messages/:message_id", func(context *gin.Context) {
		handlers.GetMessage(context, repo)
	})

	representative.PUT("/tickets/:ticket_id", func(c *gin.Context) {
		handlers.UpdateTicket(c, repo)
	})

	//statusa göre ticket döndür
	admin.GET("/tickets/", func(context *gin.Context) {
		handlers.GetTicketWith(context, repo)
	})

	//tek bir ticket görüntüleme
	authorized.GET("/tickets/:ticket_id", func(context *gin.Context) {
		handlers.GetTicket(context, repo)
	})

	//set status
	representative.PATCH("/tickets/:ticket_id", func(context *gin.Context) {
		handlers.SetTicketStatus(context, repo)
	})

	//bütün ticketleri almak için

	authorized.GET("/tickets", func(context *gin.Context) {
		handlers.GetAllTickets(context, repo)
	})
	//oluşturmak için
	authorized.POST("/tickets", func(context *gin.Context) {
		handlers.CreateTicket(context, repo)
	})

	admin.GET("/customers", func(context *gin.Context) {
		handlers.GetCustomers(context, repo)
	})

	router.POST("/customers", func(context *gin.Context) {
		handlers.Register(context, repo)
	})

	representative.GET("/customers/:customer_id", func(context *gin.Context) {
		handlers.GetCustomer(context, repo)
	})

	representative.PUT("/customers/:customers_id", func(context *gin.Context) {
		handlers.UpdateCustomer(context, repo)
	})

	admin.GET("/representatives", func(context *gin.Context) {
		handlers.GetAllRepresentatives(context, repo)
	})

	admin.POST("/representatives", func(context *gin.Context) { //sadece admin temsilci oluşturabilir
		handlers.Register(context, repo)
	})

	authorized.PUT("/representatives/:representatives_id", func(context *gin.Context) {
		handlers.UpdateRepresentative(context, repo)
	})

	authorized.GET("/representatives/:representatives_id", func(context *gin.Context) {
		handlers.GetRepresentative(context, repo)
	})

	return router
}
