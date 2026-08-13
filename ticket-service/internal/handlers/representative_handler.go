package handlers

import (
	"TMS/ticket-service/internal/models"
	"TMS/ticket-service/internal/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateRepresentative(representative *models.Representative, c *gin.Context, repo *repository.Repository) {
	err := repo.CreateRepresentative(representative)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, representative)

}
func CreateRepresentativeTx(representative *models.Representative, c *gin.Context, repo *repository.Repository) {
	err := repo.CreateRepresentative(representative)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, representative)

}

func GetRepresentative(c *gin.Context, repo *repository.Repository) {
	idString := (c.Param("representatives_id"))
	id, err := strconv.Atoi(idString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is invalid" + err.Error()})
		return
	}
	representative, error := repo.GetRepresentative(id)
	if error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": error.Error()})
		return
	}
	c.JSON(200, representative)
}
func GetAllRepresentatives(c *gin.Context, repo *repository.Repository) {
	representatives, err := repo.GetAllRepresentatives()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, representatives)
}
func UpdateRepresentative(c *gin.Context, repo *repository.Repository) {
	idString := (c.Param("representatives_id"))
	id, err := strconv.Atoi(idString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is invalid" + err.Error()})
		return
	}
	var rep models.Representative
	if err := c.ShouldBind(&rep); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cant bind representatives" + err.Error()})
		return
	}
	err = repo.UpdateRepresentative(id, &rep)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, "updated")
}
