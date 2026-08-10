package handlers

import (
	"TMS/internal/models"
	"TMS/internal/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateAnswer(c *gin.Context, repo *repository.Repository) {
	var req models.CreateAnswerRequest
	id, err := strconv.Atoi(c.Param("ticket_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ticket_id is invalid" + err.Error()})
		return
	}

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "can't bind answer " + err.Error()})
		return
	}
	repID := int(c.GetFloat64("entity_id"))
	ok,err:=repo.IsRepresentativeAssigned(id, repID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ticket_id is not assigned"})
		return
	}
	answer := models.Answer{
		TicketID:   id,
		RepID:      repID,
		AnswerText: req.AnswerText,
	}
	//burada JSON dan alınan değeri değiştiriyoruz JSON daki veri manipüle edilebilir
	err = repo.CreateAnswer(&answer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "repo"})
	}
	c.JSON(201, answer)
}
func GetAnswer(c *gin.Context, repo *repository.Repository) {
	tid, err := strconv.Atoi(c.Param("ticket_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ticket_id is invalid" + err.Error()})
		return
	}
	aid, err := strconv.Atoi(c.Param("answer_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "answer_id is invalid" + err.Error()})
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
		repId:=int(c.GetFloat64("entity_id"))
		ok,err:=repo.IsRepresentativeAssigned(tid, repId)
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

	ok, err := repo.IsAnswerBelongsToTicket(aid, tid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden answer"})
		return
	}
	answer, error := repo.GetAnswer(aid)
	if error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": error.Error()})
		return
	}
	if answer == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "answer not found"})
		return
	}
	c.JSON(200, answer)
}
func UpdateAnswer(c *gin.Context, repo *repository.Repository) {
	var req models.CreateAnswerRequest
	aid, err := strconv.Atoi(c.Param("answer_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "param: answer_id is invalid"})
		return
	}
	ticketId, err := repo.GetTicketIDByAnswerID(aid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	repID := int(c.GetFloat64("entity_id"))
	ok,err:=repo.IsRepresentativeAssigned(ticketId, repID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ticket_id is not assigned"})
		return
	}
	ok, err = repo.IsAnswerMatchWithRep(aid, repID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	oldValue, err := repo.GetAnswer(aid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil { //yine aynı mantıkla var olan verileri bindiliyoruz id yi urlden alcaz
		c.JSON(http.StatusBadRequest, gin.H{"error": "can't bind answer" + err.Error()})
		return
	}
	answer := models.Answer{
		AnswerText: req.AnswerText,
	}
	err = repo.UpdateAnswer(aid, &answer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = repo.CreateActivityLog(aid, repID, models.ActionAnswerUpdated, "description", oldValue.AnswerText, req.AnswerText)
	c.JSON(http.StatusOK, answer)

}
func GetAnswers(c *gin.Context, repo *repository.Repository) {
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
		ok, err := repo.IsRepresentativeAssignedAnswer(id, repID)
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

	answers, error := repo.GetAnswers(id)
	if error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": error.Error()})
		return
	}
	c.JSON(200, answers)

}
