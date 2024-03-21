package controllers

import (
	"net/http"
	"restaurant/database"
	"restaurant/models"

	"github.com/gin-gonic/gin"
)


func AssignTable(c *gin.Context) {
	var staff models.StaffModel
	var table models.TablesModel

	staffID := c.Param("staffID")
	tableID := c.Param("tableID")

	if err := database.DB.Find(&staff, staffID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Staff not found"})
		return
	}

	if err := database.DB.Find(&table, tableID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Table not found"})
		return
	}

	staff.Staff_id = table.TableID
	if err := database.DB.Save(&staff).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, table)
}

func CreateOrder(c *gin.Context) {
	var order models.OrderModel
	// tableID := c.Param("tableID")
	// staffID := c.Param("staffID")

	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

   
	if err := database.DB.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	c.JSON(http.StatusCreated, order)
}