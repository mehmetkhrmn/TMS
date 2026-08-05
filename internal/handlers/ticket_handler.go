package handlers

import (
	"TMS/internal/models"
	"TMS/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateTicket(ticket *models.Ticket, c *gin.Context, repo *repository.Repository) { //repositorydeki structı gönderdik bağlantı kurudk db ile

	//database ekliyoruz
	if err := repo.CreateTicket(ticket); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cant create ticket" + err.Error()})
	}
	c.JSON(http.StatusCreated, ticket) //burada da create ticket fonksyonu içide yazdığımız returning ile güncellenmiş json var

}

func GetAllTickets(c *gin.Context, repo *repository.Repository) {
	tickets, error := repo.GetAllTickets()
	if error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": error})
		return
	}
	c.JSON(http.StatusOK, tickets)
}

func GetTicketWith(status string, c *gin.Context, repo *repository.Repository) {

	tickets, error := repo.GetAllWithStatus(status)
	if error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": error.Error()})
		return
	}
	c.JSON(http.StatusOK, tickets)

}

func GetTicket(id int, c *gin.Context, repo *repository.Repository) {
	ticket, error := repo.GetTicket(id)

	if error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": error.Error()})
		return
	}
	if ticket == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ticket not found"})
		return
	}
	c.JSON(http.StatusOK, ticket)
}

func SetTicketStatus(id int, status string, c *gin.Context, repo *repository.Repository) {
	err := repo.SetStatus(id, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"id": id, "status": status})
}

func UpdateTicket(id int, ticket *models.Ticket, c *gin.Context, repo *repository.Repository) {

	err := repo.UpdateTicket(id, ticket)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "repo"})
	}
	c.JSON(200, ticket)
}
