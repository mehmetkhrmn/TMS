package handlers

import (
	"TMS/ticket-service/internal/messaging"
	"TMS/ticket-service/internal/models"
	"TMS/ticket-service/internal/repository"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

func CreateMessage(c *gin.Context, repo *repository.Repository, rabbit *messaging.RabbitMQ) {
	var req models.RequestTicketMessage
	id, err := strconv.Atoi(c.Param("ticket_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ticket_id is invalid" + err.Error()})
		return
	}

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "can't bind message " + err.Error()})
		return
	}
	userID := int(c.GetFloat64("user_id"))
	role := c.GetString("role")
	switch role {
	case "representative":
		repID := int(c.GetFloat64("entity_id"))
		ok, err := repo.IsRepresentativeAssigned(id, repID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ticket_id is not assigned"})
			return
		}

	case "customer":
		cusID := int(c.GetFloat64("entity_id"))
		ok, err := repo.IsTicketOwner(id, cusID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ticket_id is not assigned"})
			return
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role"})
		return
	}

	message := models.TicketMessage{
		TicketID: id,
		UserID:   userID,
		Message:  req.Message,
	}
	//burada JSON dan alınan değeri değiştiriyoruz JSON daki veri manipüle edilebilir
	tx, err := repo.Db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer func() {
		_ = tx.Rollback()

	}()

	err = repo.CreateMessage(tx, &message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "repo"})
		return
	}

	request := models.NotificationRequest{}
	eventID := uuid.New().String()
	skipNotif := false
	switch role {
	case "representative":
		userID2, err := repo.GetCustomerUserIDByTicketID(id) //alıcı olacak
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		err = repo.CreateActivityLog(tx, id, userID, models.ActionMessageCreated, "message", "", "")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		request = models.NotificationRequest{
			TicketID:        id,
			ActorUserID:     userID,
			RecipientUserID: userID2,
			Type:            models.ActionMessageCreated,
			Message:         "Message created",
			OccurredAt:      time.Now(),
			EventID:         eventID,
		}
	case "customer":
		userID2, err := repo.GetRepresentativeUserIDByTicketID(id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				skipNotif = true
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "server"})
				return
			}
		}
		err = repo.CreateActivityLog(tx, id, userID, models.ActionMessageCreated, "message", "", "")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if !skipNotif {
			request = models.NotificationRequest{
				TicketID:        id,
				ActorUserID:     userID,
				RecipientUserID: userID2,
				Type:            models.ActionMessageCreated,
				Message:         "Message created",
				OccurredAt:      time.Now(),
				EventID:         eventID,
			}
		}

	}
	err = tx.Commit()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if skipNotif {
		c.JSON(201, message)
		return
	}
	err = rabbit.PublishNotification(request)
	if err != nil {
		slog.Error("Notification can't publish", "error", err)
	}
	slog.Info("Publishing notification",
		"exchange", os.Getenv("RABBITMQ_EXCHANGE"),
		"routing_key", os.Getenv("RABBITMQ_ROUTING_KEY"),
		"queue", os.Getenv("RABBITMQ_QUEUE"),
	)
	c.JSON(201, message)
}
func GetMessage(c *gin.Context, repo *repository.Repository) {
	tid, err := strconv.Atoi(c.Param("ticket_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ticket_id is invalid" + err.Error()})
		return
	}
	mid, err := strconv.Atoi(c.Param("message_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "message_id is invalid" + err.Error()})
		return
	}

	role := c.GetString("role")
	switch role {
	case "customer":
		custId := int(c.GetFloat64("entity_id"))
		ok, err := repo.IsTicketOwner(tid, custId)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden ticket"})
			return
		}
	case "representative":
		repId := int(c.GetFloat64("entity_id"))
		ok, err := repo.IsRepresentativeAssigned(tid, repId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ticket_id is not assigned"})
			return
		}
	case "admin":
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role"})
		return
	}

	ok, err := repo.IsMessageBelongsToTicket(mid, tid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden message"})
		return
	}
	message, err := repo.GetMessage(mid)
	fmt.Print(message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if message == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "message not found"})
		return
	}
	c.JSON(200, message)
}

// todo : eski rest fonsiyonalrını MQ ile uyumlu oalcak şekilde değiştir
func UpdateMessage(c *gin.Context, repo *repository.Repository, rabbit *messaging.RabbitMQ) {
	var req models.UpdateMessageRequest
	tid, err := strconv.Atoi(c.Param("ticket_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ticket_id is invalid" + err.Error()})
		return
	}
	mid, err := strconv.Atoi(c.Param("message_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "param: message_id is invalid"})
		return
	}

	ok, err := repo.IsMessageBelongsToTicket(mid, tid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden message"})
		return
	}
	userID := int(c.GetFloat64("user_id"))
	role := c.GetString("role")
	eventID := uuid.New().String()
	switch role {
	case "representative":
		repID := int(c.GetFloat64("entity_id"))
		ok, err := repo.IsRepresentativeAssigned(tid, repID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ticket_id is not assigned"})
			return
		}
		ok, err = repo.IsMessageMatchWithUser(mid, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden ticket"})
			return
		}
	case "customer":
		cusID := int(c.GetFloat64("entity_id"))
		ok, err := repo.IsTicketOwner(tid, cusID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden ticket"})
			return
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role"})
		return

	}
	oldValue, err := repo.GetMessage(mid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil { //yine aynı mantıkla var olan verileri bindiliyoruz id yi urlden alcaz
		c.JSON(http.StatusBadRequest, gin.H{"error": "can't bind message" + err.Error()})
		return
	}
	message := models.UpdateMessageRequest{
		Message: req.Message,
	}
	tx, err := repo.Db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer func() {
		_ = tx.Rollback()

	}()

	err = repo.UpdateMessage(tx, mid, &message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = repo.CreateActivityLog(tx, tid, userID, models.ActionMessageUpdated, "description", oldValue.Message, req.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = tx.Commit()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	request := models.NotificationRequest{
		TicketID:    tid,
		ActorUserID: userID,
		Type:        models.ActionMessageUpdated,
		Message:     "Message updated",
		OccurredAt:  time.Now(),
		EventID:     eventID,
	}
	skipNotif := false
	switch role {
	case "representative":
		// Bildirim customer'a gidecek
		userID2, err := repo.GetCustomerUserIDByTicketID(tid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		request.RecipientUserID = userID2

	case "customer":
		// Bildirim representative'a gidecek
		userID2, err := repo.GetRepresentativeUserIDByTicketID(tid)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				skipNotif = true
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		} else {
			request.RecipientUserID = userID2
		}
	}
	if !skipNotif {
		err = rabbit.PublishNotification(request)
		if err != nil {
			slog.Error("Notification can't publish", "error", err)
		}
	}
	c.JSON(http.StatusCreated, gin.H{"message": "updated"})
}
func GetMessages(c *gin.Context, repo *repository.Repository) {
	id, err := strconv.Atoi(c.Param("ticket_id"))
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
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized ticket"})
			return
		}
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ticket_id is invalid" + err.Error()})
			return
		}
	case "representative":
		repID := int(c.GetFloat64("entity_id"))
		ok, err := repo.IsRepresentativeAssigned(id, repID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized ticket"})
			return
		}
	case "admin":
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role"})
		return
	}

	messages, err := repo.GetMessages(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, messages)

}
