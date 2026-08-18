package handlers

import (
	"TMS/ticket-service/internal/messaging"
	"TMS/ticket-service/internal/models"
	"TMS/ticket-service/internal/notification"
	"TMS/ticket-service/internal/repository"
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func CreateTicket(c *gin.Context, repo *repository.Repository, client *notification.Client, rabbit *messaging.RabbitMQ) { //repositorydeki structı gönderdik bağlantı kurudk db ile
	var req models.CreateTicketRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cant bind ticket" + err.Error()})
		return
	}
	custId := int(c.GetFloat64("entity_id"))
	userId := int(c.GetFloat64("user_id"))
	ticket := models.Ticket{
		Subject:     req.Subject,
		Description: req.Description,
		CustomerID:  custId,
		Category:    req.Category,
	}
	tx, error := repo.Db.Begin()
	if error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": error.Error()})
		return
	}
	defer func() {
		_= tx.Rollback()

	}()
	//database ekliyoruz
	if err := repo.CreateTicket(tx, &ticket); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cant create ticket" + err.Error()})
		return
	}
	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err := client.CreateNotification(notification.NotificationRequest{
		TicketID: ticket.ID,
		UserID:   userId,
		Type:     "ticket_created",
		Message:  "Your ticket has been created",
	})
	if err != nil {
		slog.Error("Notification can't sent")
	}
	c.JSON(http.StatusCreated, ticket) //burada da create ticket fonksyonu içide yazdığımız returning ile güncellenmiş json var

}

func GetAllTickets(c *gin.Context, repo *repository.Repository) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page"})
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 || limit > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
		return
	}

	offset := (page - 1) * limit

	role := c.GetString("role")

	switch role {
	case "customer":
		custID := int(c.GetFloat64("entity_id"))

		tickets, err := repo.GetCustomerTickets(custID, limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "cant get tickets"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data":  tickets,
			"page":  page,
			"limit": limit,
		})

	case "representative":
		repID := int(c.GetFloat64("entity_id"))

		tickets, err := repo.GetRepresentativeTickets(repID, limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "cant get tickets"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data":  tickets,
			"page":  page,
			"limit": limit,
		})

	case "admin":
		tickets, err := repo.GetAllTickets(limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "cant get tickets"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data":  tickets,
			"page":  page,
			"limit": limit,
		})

	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cant get tickets"})
		return
	}
}

func GetTicketWith(c *gin.Context, repo *repository.Repository) {
	status := c.Query("status")
	category := c.Query("category")

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page"})
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 || limit > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
		return
	}

	offset := (page - 1) * limit

	tickets, err := repo.GetAllWith(
		status,
		category,
		limit,
		offset,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  tickets,
		"page":  page,
		"limit": limit,
	})
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
		userID := int(c.GetFloat64("user_id"))
		ok, err := repo.IsRepresentativeAssigned(ticketID, userID)
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
	switch role {

	case "customer":

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
	case "representative":
		representativeId := int(c.GetFloat64("entity_id"))
		ok, err := repo.IsRepresentativeAssigned(id, representativeId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden ticket"})
			return
		}
	case "admin":
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "ticket is invalid"})
		return
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

func SetTicketStatus(c *gin.Context, repo *repository.Repository, client *notification.Client, rabbit *messaging.RabbitMQ) {
	idString := c.Param("ticket_id")
	statusString := c.Query("status") //gin iki tane parametre almıyo o yüzden query den alıyoz
	id, err := strconv.Atoi(idString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ticket_id is invalid" + err.Error()})
		return
	}
	status := statusString
	repId := int(c.GetFloat64("entity_id"))
	userId := int(c.GetFloat64("user_id"))

	ok, err := repo.IsRepresentativeAssigned(id, repId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden ticket"})
		return
	}
	oldTicket, err := repo.GetTicket(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	tx, err := repo.Db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = repo.SetStatus(tx, id, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = repo.CreateActivityLog(tx, id, userId, models.ActionStatusChanged, "status", oldTicket.Status, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err= tx.Commit()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = client.CreateNotification(notification.NotificationRequest{
		TicketID: id,
		UserID:   userId,
		Type:     "ticket_updated",
		Message:  "Ticket status updated to -> " + statusString,
	})
	if err != nil {
		slog.Error("Notification can't sent")
	}
	c.JSON(200, gin.H{"id": id, "status": status})
}

func UpdateTicket(c *gin.Context, repo *repository.Repository, client *notification.Client, rabbit *messaging.RabbitMQ) {
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
	var req models.TicketUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil { //burada bindliyoruz biz gönderilmeyen veriler boş kalcak bu sayede
		c.JSON(http.StatusBadRequest, gin.H{"error": "can't bind ticket"})
		return
	}
	ticket := models.Ticket{
		ID:          oldTicket.ID,
		Status:      req.Status,
		Description: req.Description,
		Subject:     req.Subject,
		Category:    req.Category,
		CustomerID:  oldTicket.CustomerID,
	}
	tx, err := repo.Db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = repo.UpdateTicket(tx, id, &ticket)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "repo"})
		return
	}
	userId := int(c.GetFloat64("user_id"))
	err = repo.CreateActivityLog(tx, id, userId, models.ActionTicketUpdated, "description", oldTicket.Description, oldTicket.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "repo"})
		return
	}
	err= tx.Commit()
	if err!= nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = client.CreateNotification(notification.NotificationRequest{
		TicketID: id,
		UserID:   userId,
		Type:     "ticket_updated",
		Message:  "Ticket updated",
	})
	if err != nil {
		slog.Error("Notification can't sent")
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "representative id is invalid"})
		return
	}
	ticketID, err := strconv.Atoi(c.Param("ticket_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ticket_id is invalid"})
		return
	}
	tx, err := repo.Db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = repo.AssignRepresentative(tx, ticketID, repID)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			c.JSON(http.StatusConflict, gin.H{
				"error": "representative is already assigned to this ticket",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ticket or representative not found"})
		return

	}
	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"assigned": true})
}
func UnassignRepresentative(c *gin.Context, repo *repository.Repository) {
	repID, err := strconv.Atoi(c.Param("representative_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "representative id is invalid"})
		return
	}
	ticketID, err := strconv.Atoi(c.Param("ticket_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ticket_id is invalid"})
		return
	}
	tx, err := repo.Db.Begin()
	defer func(){
		_=tx.Rollback()

	}()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = repo.UnAssignRepresentative(tx, ticketID, repID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "representative is not assigned to this ticket",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ticket or representative not found"})
		return
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"unassigned": true})

}
