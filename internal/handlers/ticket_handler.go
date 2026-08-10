package handlers

import (
	"TMS/internal/models"
	"TMS/internal/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateTicket(c *gin.Context, repo *repository.Repository) { //repositorydeki structı gönderdik bağlantı kurudk db ile
	var req models.CreateTicketRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cant bind ticket" + err.Error()})
		return
	}
	custId := int(c.GetFloat64("entity_id"))
	ticket := models.Ticket{
		Subject:     req.Subject,
		Description: req.Description,
		CustomerID:  custId,
		Status:      "open",
	}
	//database ekliyoruz
	if err := repo.CreateTicket(&ticket); err != nil {
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

func GetTicketWith(c *gin.Context, repo *repository.Repository) {
	status := c.Query("ticket_status")
	tickets, error := repo.GetAllWithStatus(status)
	if error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": error.Error()})
		return
	}
	c.JSON(http.StatusOK, tickets)

}
func GetTicketHistory(c *gin.Context, repo *repository.Repository) {
	ticketID, err := strconv.Atoi(c.Param("ticket_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ticket is invalid" + err.Error()})
		return
	}
	role := c.GetString("role")
	switch role {
	case "customer":
		custId := int(c.GetFloat64("entity_id"))
		ok, err := repo.IsTicketOwner(ticketID, custId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "ticket is not owner"})
			return
		}
		answers, err := repo.GetTicketHistory(ticketID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, answers)

	case "representative":
		representativeID := int(c.GetFloat64("entity_id"))
		ok, err := repo.IsRepresentativeAssignedAnswer(ticketID, representativeID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "representative is not assigned answer"})
			return
		}
		answers, err := repo.GetTicketHistory(ticketID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, answers)
	case "admin":
		answers, err := repo.GetTicketHistory(ticketID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, answers)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "ticket is invalid"})
		return
	}

}

func GetTicket(c *gin.Context, repo *repository.Repository) {
	idString := (c.Param("ticket_id")) //url den id yi aldık
	id, err := strconv.Atoi(idString)  //int ye çevir
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ticket_id is invalid" + err.Error()})
		return
	}
	role := c.GetString("role")
	if role == "customer" {

		custId := int(c.GetFloat64("entity_id"))
		ok, err := repo.IsTicketOwner(id, custId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden ticket"})
			return
		}

	}
	ticket, err := repo.GetTicket(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if ticket == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ticket not found"})
		return
	}

	c.JSON(http.StatusOK, ticket)
}

func SetTicketStatus(c *gin.Context, repo *repository.Repository) {
	idString := c.Param("ticket_id")
	statusString := c.Query("status") //gin iki tane parametre almıyo o yüzden query den alıyoz
	id, err := strconv.Atoi(idString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ticket_id is invalid" + err.Error()})
		return
	}
	status := statusString
	repId := int(c.GetFloat64("entity_id"))
	ok, err := repo.IsRepresentativeAssigned(id, repId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden ticket"})
		return
	}

	err = repo.SetStatus(id, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"id": id, "status": status})
}

func UpdateTicket(c *gin.Context, repo *repository.Repository) {
	id, err := strconv.Atoi(c.Param("ticket_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ticket_id is invalid" + err.Error()})
		return
	}
	repID := int(c.GetFloat64("entity_id"))
	ok, err := repo.IsRepresentativeAssigned(id, repID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden ticket"})
		return
	}
	oldTicket, error := repo.GetTicket(id)
	if error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": error.Error()})
		return
	}
	var ticket models.Ticket
	if err := c.ShouldBindJSON(&ticket); err != nil { //burada bindliyoruz biz gönderilmeyen veriler boş kalcak bu sayede
		c.JSON(http.StatusBadRequest, gin.H{"error": "can't bind ticket" + err.Error()})
		return
	}
	err = repo.UpdateTicket(id, &ticket)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "repo"})
	}
	err = repo.CreateActivityLog(id, repID, models.ActionTicketUpdated, "description", oldTicket.Description, oldTicket.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "repo"})
		return
	}
	c.JSON(200, ticket)
}

func GetAdminTickets(ticket_id int, c *gin.Context, repo *repository.Repository) {
	tickets, err := repo.GetAdminTicket(ticket_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "repo"})
		return
	}
	c.JSON(http.StatusOK, tickets)
}

func GetRepresentativeTickets(customer_id int, representative_id int, c *gin.Context, repo *repository.Repository) {
	tickets, err := repo.GetRepresentativeTickets(customer_id, representative_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "repo"})
		return
	}

	c.JSON(http.StatusOK, tickets)
}

func GetCustomerTicket(ticket_id int, customer_id int, c *gin.Context, repo *repository.Repository) {
	ticket, err := repo.GetCustomerTicket(ticket_id, customer_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "repo"})
		return
	}
	c.JSON(http.StatusOK, ticket)
}
func AssignRepresentative(c *gin.Context, repo *repository.Repository) {
	repID, err := strconv.Atoi(c.Param("representative_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "representative id is invalid" + err.Error()})
		return
	}
	ticketID, err := strconv.Atoi(c.Param("ticket_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ticket_id is invalid" + err.Error()})
		return
	}

	err = repo.AssignRepresentative(ticketID, repID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "repo"})
		return
	}
	c.JSON(200, gin.H{"assigned": true})
}
