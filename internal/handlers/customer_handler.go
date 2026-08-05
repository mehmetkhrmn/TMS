package handlers

import (
	"TMS/internal/models"
	"TMS/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateCustomer(customer *models.Customer, c *gin.Context, repo *repository.Repository) {
	err := repo.CreateCustomer(customer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(201, customer)

}
func CreateCustomerTx(customer *models.Customer, c *gin.Context, repo *repository.Repository) {
	err := repo.CreateCustomer(customer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(201, customer)

}
func GetCustomers(c *gin.Context, repo *repository.Repository) {
	customers, error := repo.GetCustomers()
	if error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": error.Error()})
		return
	}
	c.JSON(200, customers)
}
func GetCustomer(id int, c *gin.Context, repo *repository.Repository) {
	customer, error := repo.GetCustomer(id)
	if error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": error.Error()})
		return
	}
	c.JSON(200, customer)
}
func UpdateCustomer(id int, customer *models.Customer, c *gin.Context, repo *repository.Repository) {
	_, err := repo.UpdateCustomer(id, customer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, "updated")
}
