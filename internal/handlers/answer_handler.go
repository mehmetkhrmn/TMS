package handlers

import (
	"TMS/internal/models"
	"TMS/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateAnswer(answer *models.Answer, c *gin.Context, repo *repository.Repository) {
	err := repo.CreateAnswer(answer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "repo"})
	}
	c.JSON(201, answer)
}
func GetAnswer(ticketId int, answerId int, c *gin.Context, repo *repository.Repository) {
	answer, error := repo.GetAnswer(ticketId, answerId)
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
func UpdateAnswer(id int, ticket *models.Answer, c *gin.Context, repo *repository.Repository) {
	err := repo.UpdateAnswer(id, ticket)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, ticket)

}
func GetAnswers(id int, c *gin.Context, repo *repository.Repository) {
	answers, error := repo.GetAnswers(id)
	if error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": error.Error()})
		return
	}

	c.JSON(200, answers)
}
