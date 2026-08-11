package handlers

import (
	"TMS/internal/models"
	"TMS/internal/repository"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateMessage(c *gin.Context, repo *repository.Repository) {
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
	err = repo.CreateMessage(&message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "repo"})
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
func UpdateMessage(c *gin.Context, repo *repository.Repository) {
	var req models.UpdateMessageRequest
	tid, err := strconv.Atoi(c.Param("ticket_id"))
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
	err = repo.UpdateMessage(mid, &message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = repo.CreateActivityLog(mid, userID, models.ActionMessageUpdated, "description", oldValue.Message, req.Message)
	c.JSON(http.StatusOK, message)

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
