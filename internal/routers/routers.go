package routers

import (
	"TMS/internal/models"
	"TMS/internal/repository"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter(db *sql.DB) *gin.Engine {

	router := gin.Default()
	gin.SetMode(gin.DebugMode)
	repo := repository.NewRepository(db)

	//tickets
	//görüntülemek için
	router.GET("/tickets")
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
	c.JSON(201, ticket) //buradada create ticket fonksyonu içide yazdığımız returning ile güncellenmiş json var

}
