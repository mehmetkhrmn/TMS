package routers

import (
	"TMS/internal/handlers"
	"TMS/internal/middleware"
	"TMS/internal/models"
	"TMS/internal/repository"

	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func SetupRouter(db *sql.DB) *gin.Engine {

	router := gin.Default()
	gin.SetMode(gin.DebugMode)
	repo := repository.NewRepository(db)
	authorized := router.Group("/") //route için alt grup oluşturuyorum
	authorized.Use(middleware.AuthMiddleware())

	router.POST("/login", func(context *gin.Context) {
		handlers.Login(context, repo)
	})
	authorized.PUT("/answers/:answer_id", func(context *gin.Context) {
		id, err := strconv.Atoi(context.Param("answer_id"))
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": "param: answer_id is invalid"})
			return
		}
		var answer models.Answer
		if err := context.ShouldBind(&answer); err != nil { //yine aynı mantıkla var olan verileri bindiliyoruz id yi urlden alcaz
			context.JSON(http.StatusBadRequest, gin.H{"error": "can't bind answer" + err.Error()})
			return
		}
		handlers.UpdateAnswer(id, &answer, context, repo)
	})
	authorized.POST("/tickets/:ticket_id/answers", func(context *gin.Context) {
		id, err := strconv.Atoi(context.Param("ticket_id"))
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": "ticket_id is invalid" + err.Error()})
			return
		}
		var answer models.Answer
		if err := context.ShouldBind(&answer); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": "can't bind answer" + err.Error()})
			return
		}
		answer.TicketID = id //burada JSON dan alınan değeri değiştiriyoruz JSON daki veri manipüle edilebilir
		handlers.CreateAnswer(&answer, context, repo)
	})

	authorized.GET("/tickets/:ticket_id/answers", func(context *gin.Context) {
		id, err := strconv.Atoi(context.Param("ticket_id"))
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": "ticket_id is invalid" + err.Error()})
			return
		}
		handlers.GetAnswers(id, context, repo)
	})
	authorized.GET("/tickets/:ticket_id/answers/:answer_id", func(context *gin.Context) {
		tid, err := strconv.Atoi(context.Param("ticket_id"))

		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": "ticket_id is invalid" + err.Error()})
			return
		}
		aid, err := strconv.Atoi(context.Param("answer_id"))
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": "answer_id is invalid" + err.Error()})
			return
		}
		handlers.GetAnswer(tid, aid, context, repo)
	})
	authorized.PUT("/tickets/:ticket_id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("ticket_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ticket_id is invalid" + err.Error()})
			return
		}
		var ticket models.Ticket
		if err := c.ShouldBindJSON(&ticket); err != nil { //burada bindliyoruz biz gönderilmeyen veriler boş kalcak bu sayede
			c.JSON(http.StatusBadRequest, gin.H{"error": "can't bind ticket" + err.Error()})
			return
		}
		handlers.UpdateTicket(id, &ticket, c, repo)
	})
	//statusa göre ticket döndür
	authorized.GET("tickets/", func(context *gin.Context) {
		status := context.Query("ticket_status")
		role := context.GetString("role")
		if role != "admin" && role != "representative" {
			context.JSON(http.StatusUnauthorized, gin.H{"error": "role is required"})
			return
		}
		switch status {
		case "open", "closed", "in_progress", "resolved":
			handlers.GetTicketWith(status, context, repo)

		case "":
			handlers.GetAllTickets(context, repo)
		default:
			context.JSON(http.StatusBadRequest, gin.H{"error": "ticket_status is unknown -> " + status})
			return
		}
	})

	//tek bir ticket görüntüleme
	authorized.GET("/tickets/:ticket_id", func(context *gin.Context) {
		idString := (context.Param("ticket_id")) //url den id yi aldık
		id, err := strconv.Atoi(idString)        //int ye çevir
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": "ticket_id is invalid" + err.Error()})
			return
		}
		handlers.GetTicket(id, context, repo)
	})
	//set status
	authorized.PATCH("tickets/:ticket_id", func(context *gin.Context) {
		idString := context.Param("ticket_id")
		statusString := context.Query("status") //gin iki tane parametre almıyo o yüzden query den alıyoz
		id, err := strconv.Atoi(idString)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": "ticket_id is invalid" + err.Error()})
			return
		}
		status := statusString

		handlers.SetTicketStatus(id, status, context, repo)

	})
	//bütün ticketleri almak için
	authorized.GET("/tickets", func(context *gin.Context) {
		handlers.GetAllTickets(context, repo)
	})
	//oluşturmak için
	authorized.POST("/tickets", func(context *gin.Context) {
		var ticket models.Ticket
		if err := context.ShouldBind(&ticket); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": "cant bind ticket" + err.Error()})
			return
		}
		handlers.CreateTicket(&ticket, context, repo)
	})
	authorized.GET("/customers", func(context *gin.Context) {
		handlers.GetCustomers(context, repo)
	})
	router.POST("/customers", func(context *gin.Context) {

		handlers.Register(context, repo)

	})
	authorized.GET("/customers/:customer_id", func(context *gin.Context) {
		idString := (context.Param("customer_id"))
		id, err := strconv.Atoi(idString)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": "customer_id is invalid" + err.Error()})
			return
		}
		handlers.GetCustomer(id, context, repo)

	})
	authorized.PUT("/customers/:customers_id", func(context *gin.Context) {
		idString := (context.Param("customers_id"))
		id, err := strconv.Atoi(idString)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": "customer_id is invalid" + err.Error()})
			return
		}
		var customer models.Customer
		if err := context.ShouldBind(&customer); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": "cant bind customer " + err.Error()})
			return
		}
		handlers.UpdateCustomer(id, &customer, context, repo)
	})
	authorized.GET("/representatives", func(context *gin.Context) {
		handlers.GetAllRepresentatives(context, repo)
	})
	authorized.POST("/representatives", func(context *gin.Context) { //sadece admin temsilci oluşturabilir
		handlers.Register(context, repo)
	})
	authorized.PUT("/representatives/:representatives_id", func(context *gin.Context) {
		idString := (context.Param("representatives_id"))
		id, err := strconv.Atoi(idString)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": "id is invalid" + err.Error()})
			return
		}
		var rep models.Representative
		if err := context.ShouldBind(&rep); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": "cant bind representatives" + err.Error()})
			return
		}
		handlers.UpdateRepresentative(id, &rep, context, repo)
	})
	authorized.GET("/representatives/:representatives_id", func(context *gin.Context) {
		idString := (context.Param("representatives_id"))
		id, err := strconv.Atoi(idString)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": "id is invalid" + err.Error()})
			return
		}
		handlers.GetRepresentative(id, context, repo)
	})
	return router
}
