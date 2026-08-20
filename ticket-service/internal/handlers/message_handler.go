package handlers

import (
	"TMS/ticket-service/internal/messaging"
	"TMS/ticket-service/internal/models"
	"TMS/ticket-service/internal/notification"
	"TMS/ticket-service/internal/repository"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
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
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

	}
	err = repo.CreateMessage(tx, &message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "repo"})
		return
	}

	var mesType string
	if role == "representative" {
		mesType = "message_replied"
	} else {
		mesType = "message_updated"
	}

	request := notification.NotificationRequest{}
	eventID := uuid.New().String()
	switch role {
	case "representative":
		userID2, err := repo.GetCustomerUserIDByTicketID(id) //alıcı olacak
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		err = repo.CreateActivityLog(tx, id, userID, "message_created", "message", "", "")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		request = notification.NotificationRequest{
			TicketID:        id,
			ActorUserID:     userID,
			RecipientUserID: userID2,
			Type:            mesType,
			Message:         "Message created",
			OccurredAt:      time.Now(),
			EventID:         eventID,
		}
	case "customer":
		userID2, err := repo.GetRepresentativeUserIDByTicketID(id)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusBadRequest, gin.H{"error": "representative not assigned"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "server"})
			return
		}
		err = repo.CreateActivityLog(tx, id, userID, "message_created", "message", "", "")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		request = notification.NotificationRequest{
			TicketID:        id,
			ActorUserID:     userID,
			RecipientUserID: userID2,
			Type:            mesType,
			Message:         "Message created",
			OccurredAt:      time.Now(),
			EventID:         eventID,
		}
	}
	err = tx.Commit()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = rabbit.Publish(request)
	if err != nil {
		slog.Error("Notification can't publish", "error", err)
	}
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
	message, error := repo.GetMessage(mid)
	fmt.Print(message)
	if error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": error.Error()})
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

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
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

	ok, err = repo.IsMessageMatchWithUser(mid, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
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
	defer func() {
		_ = tx.Rollback()

	}()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
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
	request := notification.NotificationRequest{
		TicketID:    tid,
		ActorUserID: userID,
		Type:        "message_updated",
		Message:     "Message updated",
		OccurredAt:  time.Now(),
		EventID:     eventID,
	}

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
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		request.RecipientUserID = userID2
	}

	err = rabbit.Publish(request)
	if err != nil {
		slog.Error("Notification can't publish", "error", err)
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

	messages, error := repo.GetMessages(id)
	if error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": error.Error()})
		return
	}
	c.JSON(200, messages)

}
