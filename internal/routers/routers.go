package routers

import (
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
	router.LoadHTMLGlob("internal/templates/*")
	repo := repository.NewRepository(db)

	//tickets
	//update

	router.PUT("/answers/:answer_id", func(context *gin.Context) {
		id, err := strconv.Atoi(context.Param("answer_id"))
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": "answer_id"})
			return
		}
		var answer models.Answer
		if err := context.ShouldBind(&answer); err != nil { //yine aynı mantıkla var olan verileri bindiliyoruz id yi urlden alcaz
			context.JSON(http.StatusBadRequest, gin.H{"error": "bind"})
			return
		}
		updateAnswer(id, &answer, context, repo)
	})
	router.POST("/tickets/:ticket_id/answers", func(context *gin.Context) {
		id, err := strconv.Atoi(context.Param("ticket_id"))
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": "ticket_id"})
			return
		}
		var answer models.Answer
		if err := context.ShouldBind(&answer); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": "bind"})
			return
		}
		answer.TicketID = id //burada JSON dan alınan değeri değiştiriyoruz JSON daki veri manipüle edilebilir
		createAnswer(&answer, context, repo)
	})

	router.GET("/tickets/:ticket_id/answers", func(context *gin.Context) {
		id, err := strconv.Atoi(context.Param("ticket_id"))
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": "ticket_id"})
			return
		}
		getAnswers(id, context, repo)
	})
	router.GET("/tickets/:ticket_id/answers/:answer_id", func(context *gin.Context) {
		tid, err := strconv.Atoi(context.Param("ticket_id"))

		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": "ticket_id"})
			return
		}
		aid, err := strconv.Atoi(context.Param("answer_id"))
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": "answer_id"})
			return
		}
		getAnswer(tid, aid, context, repo)
	})
	router.PUT("/tickets/:ticket_id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("ticket_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ticket_id"})
			return
		}
		var ticket models.Ticket
		if err := c.ShouldBindJSON(&ticket); err != nil { //burada bindliyoruz biz gönderilmeyen veriler boş kalcak bu sayede
			c.JSON(http.StatusBadRequest, gin.H{"error": "bind"})
			return
		}
		updateTicket(id, &ticket, c, repo)
	})
	//statusa göre ticket döndür
	router.GET("tickets/", func(context *gin.Context) {
		status := context.Query("ticket_status")

		switch status {
		case "open", "closed", "in_progress", "resolved":
			getTicketWith(status, context, repo)

		case "":
			getAllTickets(context, repo)
		default:
			context.JSON(http.StatusBadRequest, gin.H{"error": "ticket_status is unknown"})
			return
		}
	})

	//tek bir ticket görüntüleme
	router.GET("/tickets/:ticket_id", func(context *gin.Context) {
		idString := (context.Param("ticket_id")) //url den id yi aldık
		id, err := strconv.Atoi(idString)        //int ye çevir
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		getTicket(id, context, repo)
	})
	//set status
	router.PATCH("tickets/:ticket_id", func(context *gin.Context) {
		idString := context.Param("ticket_id")
		statusString := context.Query("status") //gin iki tane parametre almıyo o yüzden query den alıyoz
		id, err := strconv.Atoi(idString)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		status := statusString
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		setTicketStatus(id, status, context, repo)

	})
	//bütün ticketleri almak için
	router.GET("/tickets", func(context *gin.Context) {
		getAllTickets(context, repo)
	})
	//oluşturmak için
	router.POST("/tickets", func(context *gin.Context) {
		createTicket(context, repo)
	})

	return router
}
func createTicket(c *gin.Context, repo *repository.Repository) { //repositorydeki structı gönderdik bağlantı kurudk db ile
	var ticket models.Ticket
	//gelen json datasını tickete göre formatladık
	if err := c.BindJSON(&ticket); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//database ekliyoruz
	if err := repo.CreateTicket(&ticket); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.JSON(201, ticket) //burada da create ticket fonksyonu içide yazdığımız returning ile güncellenmiş json var

}

//todo:bütün ticketler
func getAllTickets(c *gin.Context, repo *repository.Repository) {
	tickets, error := repo.GetAllTickets()
	if error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": error})
		return
	}
	c.JSON(200, tickets)
}

//todo:açık olan ticketleri  getir
func getTicketWith(status string, c *gin.Context, repo *repository.Repository) {

	tickets, error := repo.GetAllWithStatus(status)
	if error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error()})
		return
	}
	c.JSON(http.StatusOK, tickets)

}

//todo:tek bir ticketi görüntüle
func getTicket(id int, c *gin.Context, repo *repository.Repository) {
	ticket, error := repo.GetTicket(id)
	if error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error()})
		return
	}
	c.JSON(http.StatusOK, ticket)
}

// todo ticketi editle
func setTicketStatus(id int, status string, c *gin.Context, repo *repository.Repository) {
	err := repo.SetStatus(id, status)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"id": id, "status": status})
}

func updateTicket(id int, ticket *models.Ticket, c *gin.Context, repo *repository.Repository) {

	err := repo.UpdateTicket(id, ticket)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "repo"})
	}
	c.JSON(201, "updated")
}

// get update create
func createAnswer(answer *models.Answer, c *gin.Context, repo *repository.Repository) {
	err := repo.CreateAnswer(answer)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "repo"})
	}
	c.JSON(201, answer)
}
func getAnswer(ticketId int, answerId int, c *gin.Context, repo *repository.Repository) {
	answer, error := repo.GetAnswer(ticketId, answerId)
	if error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error()})
		return
	}
	c.JSON(200, answer)
}
func updateAnswer(id int, ticket *models.Answer, c *gin.Context, repo *repository.Repository) {
	err := repo.UpdateAnswer(id, ticket)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "repo"})
	}
	c.JSON(201, "updated")

}
func getAnswers(id int, c *gin.Context, repo *repository.Repository) {
	answers, error := repo.GetAnswers(id)
	if error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error()})
		return
	}
	c.JSON(200, answers)
}
