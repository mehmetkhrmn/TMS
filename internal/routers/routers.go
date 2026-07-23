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
	router.PUT("/tickets/update/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id"})
			return
		}
		var ticket models.Ticket
		if err := c.ShouldBindJSON(&ticket); err != nil { //burada bindliyoruz biz gönderilmeyen veriler boş kalcak bu sayede
			c.JSON(http.StatusBadRequest, gin.H{"error": "bind"})
			return
		}
		updateTicket(id, &ticket, c, repo)
	})
	//açık ticket görüntüle
	router.GET("tickets/open", func(context *gin.Context) {
		getOpenTickets(context, repo)

	})

	//tek bir ticket görüntüleme
	router.GET("/tickets/:id", func(context *gin.Context) {
		idString := context.Param("id")   //url den id yi aldık
		id, err := strconv.Atoi(idString) //int ye çevir
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		getTicket(id, context, repo)
	})
	//setDone
	router.POST("tickets/:id", func(context *gin.Context) {
		idString := context.Param("id")
		doneString := context.Query("done") //gin iki tane parametre almıyo o yüzden query den alıyoz
		id, err := strconv.Atoi(idString)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		done, err := strconv.ParseBool(doneString)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		setTicketDone(id, done, context, repo)

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
func getOpenTickets(c *gin.Context, repo *repository.Repository) {
	tickets, error := repo.GetAllOpen()
	if error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error()})
		return
	}
	c.JSON(http.StatusOK, tickets) //ticket structuresinde JSON olarak ta tanımlama yaptığımız için direkt ona göre parselanıyor

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
func setTicketDone(id int, isDone bool, c *gin.Context, repo *repository.Repository) {
	err := repo.SetDone(id, isDone)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, gin.H{"id": id, "done": isDone})
}

func updateTicket(id int, ticket *models.Ticket, c *gin.Context, repo *repository.Repository) {

	err := repo.Update(id, ticket)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "repo"})
	}
	c.JSON(201, "updated")
}
