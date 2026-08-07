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

	representative.PUT("/answers/:answer_id", func(context *gin.Context) {
		handlers.UpdateAnswer(context, repo)
	})

	representative.POST("/tickets/:ticket_id/answers", func(context *gin.Context) {
		handlers.CreateAnswer(context, repo)
	})

	authorized.GET("/tickets/:ticket_id/answers", func(context *gin.Context) {
		handlers.GetAnswers(context, repo)
	})

	authorized.GET("/tickets/:ticket_id/answers/:answer_id", func(context *gin.Context) {
		handlers.GetAnswer(context, repo)
	})

	representative.PUT("/tickets/:ticket_id", func(c *gin.Context) {
		handlers.UpdateTicket(c, repo)
	})

	//statusa göre ticket döndür
	authorized.GET("/tickets/", func(context *gin.Context) {
		handlers.GetTicketWith(context, repo)
	})

	//tek bir ticket görüntüleme
	authorized.GET("/tickets/:ticket_id", func(context *gin.Context) {
		handlers.GetTicket(context, repo)
	})

	//set status
	authorized.PATCH("/tickets/:ticket_id", func(context *gin.Context) {
		handlers.SetTicketStatus(context, repo)
	})

	//bütün ticketleri almak için
	admin.GET("/tickets", func(context *gin.Context) {
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
