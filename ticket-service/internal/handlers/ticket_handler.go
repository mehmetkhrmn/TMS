package handlers

import (
	"TMS/ticket-service/internal/messaging"
	"TMS/ticket-service/internal/models"
	"TMS/ticket-service/internal/repository"
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateTicket(c *gin.Context, repo *repository.Repository, rabbit *messaging.RabbitMQ) { //repositorydeki structı gönderdik bağlantı kurudk db ile
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
	tx, err := repo.Db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer func() {
		_ = tx.Rollback()

	}()
	//database ekliyoruz
	if err := repo.CreateTicket(tx, &ticket); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cant create ticket" + err.Error()})
		return
	}
	err = repo.CreateActivityLog(tx, ticket.ID, userId, models.ActionTicketCreated, "ticket", "", "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	eventID := uuid.New().String()
	err = rabbit.PublishNotification(models.NotificationRequest{
		TicketID:        ticket.ID,
		ActorUserID:     userId,
		RecipientUserID: userId,
		Type:            models.ActionTicketCreated,
		Message:         "Ticket created",
		EventID:         eventID,
		OccurredAt:      time.Now(),
	})
	if err != nil {
		slog.Error("Notification can't publish", "error", err)
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
		repId := int(c.GetFloat64("entity_id"))
		ok, err := repo.IsRepresentativeAssigned(ticketID, repId)
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
	idString := c.Param("ticket_id")  //url den id yi aldık
	id, err := strconv.Atoi(idString) //int ye çevir
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

func SetTicketStatus(c *gin.Context, repo *repository.Repository, rabbit *messaging.RabbitMQ) {
	idString := c.Param("ticket_id")
	statusString := c.Query("status") //gin iki tane parametre almıyo o yüzden query den alıyoz
	id, err := strconv.Atoi(idString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ticket_id is invalid" + err.Error()})
		return
	}
	status := statusString
	switch status {
	case "open", "in_progress", "resolved", "closed":
		// geçerli
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status"})
		return
	}
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
	defer func() {
		_ = tx.Rollback()
	}()

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
	err = tx.Commit()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	recipId, err := repo.GetCustomerUserIDByTicketID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	eventID := uuid.New().String()

	err = rabbit.PublishNotification(models.NotificationRequest{
		EventID:         eventID,
		TicketID:        id,
		ActorUserID:     userId,
		RecipientUserID: recipId,
		Type:            models.ActionTicketUpdated,
		Message:         "Ticket status updated to -> " + statusString,
		OccurredAt:      time.Now(),
	})

	if err != nil {
		slog.Error("Notification can't publish", "error", err)
	}
	c.JSON(200, gin.H{"id": id, "status": status})
}

func UpdateTicket(c *gin.Context, repo *repository.Repository, rabbit *messaging.RabbitMQ) {
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
	oldTicket, err := repo.GetTicket(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var req models.TicketUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil { //burada bindliyoruz biz gönderilmeyen veriler boş kalcak bu sayede
		c.JSON(http.StatusBadRequest, gin.H{"error": "can't bind ticket"})
		return
	}
	ticket := models.Ticket{
		ID:          oldTicket.ID,
		Description: oldTicket.Description,
		Subject:     oldTicket.Subject,
		Category:    oldTicket.Category,
		Status:      oldTicket.Status,
		CustomerID:  oldTicket.CustomerID,
	}
	if req.Description != "" {
		ticket.Description = req.Description
	}

	if req.Subject != "" {
		ticket.Subject = req.Subject
	}

	if req.Category != "" {
		ticket.Category = req.Category
	}

	tx, err := repo.Db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer func() {
		_ = tx.Rollback()
	}()

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
	err = tx.Commit()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	eventID := uuid.New().String()
	recipId, err := repo.GetCustomerUserIDByTicketID(id)
	err = rabbit.PublishNotification(models.NotificationRequest{
		TicketID:        id,
		ActorUserID:     userId,
		RecipientUserID: recipId,
		Type:            models.ActionTicketUpdated,
		Message:         "Ticket updated",
		EventID:         eventID,
		OccurredAt:      time.Now(),
	})
	if err != nil {
		slog.Error("Notification can't publish", "error", err)
	}
	c.JSON(200, ticket)
}

func AssignRepresentative(rabbit *messaging.RabbitMQ, c *gin.Context, repo *repository.Repository) {
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
	defer func() {
		_ = tx.Rollback()
	}()

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
	userId, err := repo.GetUserIDByRepID(repID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "repo"})
		return
	}
	err = repo.CreateActivityLog(tx, ticketID, userId, models.ActionTicketGranted, "assigned", "", strconv.Itoa(ticketID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	eventID := uuid.New().String()
	err = rabbit.PublishNotification(models.NotificationRequest{
		TicketID:        ticketID,
		RecipientUserID: userId,
		ActorUserID:     0,
		Type:            models.ActionTicketGranted,
		Message:         "Ticket assigned",
		EventID:         eventID,
		OccurredAt:      time.Now(),
	})
	if err != nil {
		slog.Error("Notification can't publish", "error", err)
	}
	c.JSON(200, gin.H{"assigned": true})
}
func UnassignRepresentative(rabbit *messaging.RabbitMQ, c *gin.Context, repo *repository.Repository) {
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
	defer func() {
		_ = tx.Rollback()

	}()

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
	userId, err := repo.GetUserIDByRepID(repID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "repo"})
		return
	}
	err = repo.CreateActivityLog(tx, ticketID, userId, models.ActionTicketRevoked, "unassigned", strconv.Itoa(ticketID), "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	eventID := uuid.New().String()

	notif := models.NotificationRequest{
		EventID:         eventID,
		TicketID:        ticketID,
		RecipientUserID: userId,
		Type:            models.ActionTicketRevoked,
		Message:         "Ticket revoked",
		OccurredAt:      time.Now(),
	}

	err = rabbit.PublishNotification(notif)
	if err != nil {
		slog.Error("Notification can't publish", "error", err)
	}
	c.JSON(200, gin.H{"unassigned": true})

}
