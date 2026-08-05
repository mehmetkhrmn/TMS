package handlers

import (
	"TMS/internal/models"
	"TMS/internal/repository"
	"net/http"

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

func GetRepresentative(id int, c *gin.Context, repo *repository.Repository) {
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
func UpdateRepresentative(id int, representative *models.Representative, c *gin.Context, repo *repository.Repository) {
	err := repo.UpdateRepresentative(id, representative)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, "updated")
}
