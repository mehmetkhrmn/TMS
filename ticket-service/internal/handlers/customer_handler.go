package handlers

import (
	"TMS/ticket-service/internal/models"
	"TMS/ticket-service/internal/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetCustomers(c *gin.Context, repo *repository.Repository) {
	customers, err := repo.GetCustomers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, customers)
}
func GetCustomer(c *gin.Context, repo *repository.Repository) {
	idString := (c.Param("customer_id"))
	id, err := strconv.Atoi(idString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "customer_id is invalid" + err.Error()})
		return
	}
	role := c.GetString("role")
	switch role {
	case "representative":
		repId := int(c.GetFloat64("entity_id"))
		assigned, err := repo.IsCustomerAssignedToRepresentative(id, repId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if !assigned {
			c.JSON(http.StatusForbidden, gin.H{"error": "customer is not assigned to this representative"})
			return
		}
	case "customer":
		custId := int(c.GetFloat64("entity_id"))
		if custId != id {
			c.JSON(http.StatusForbidden, gin.H{"error": "can't access this customer with this customer id"})
			return
		}
	case "admin":
	default:
		c.JSON(http.StatusForbidden, gin.H{"error": "can't access this customer with this customer id"})
		return
	}

	customer, err := repo.GetCustomer(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, customer)
}
func UpdateCustomer(c *gin.Context, repo *repository.Repository) {
	idString := (c.Param("customers_id"))
	id, err := strconv.Atoi(idString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "customer_id is invalid" + err.Error()})
		return
	}
	role := c.GetString("role")
	switch role {
	case "representative":
		repId := int(c.GetFloat64("entity_id"))
		assigned, err := repo.IsCustomerAssignedToRepresentative(id, repId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if !assigned {
			c.JSON(http.StatusForbidden, gin.H{"error": "customer is not assigned to this representative"})
			return
		}
	case "admin":
	case "customer":
		custId := int(c.GetFloat64("entity_id"))
		if custId != id {
			c.JSON(http.StatusForbidden, gin.H{"error": "can't access this customer with this customer id"})
			return
		}
	default:
		c.JSON(http.StatusForbidden, gin.H{"error": "can't access this customer with this customer id"})
		return
	}
	var customer models.Customer
	if err := c.ShouldBind(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cant bind customer " + err.Error()})
		return
	}
	_, err = repo.UpdateCustomer(id, &customer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, "updated")
}
