package controllers

import (
	//"fmt"
	"net/http"
	"restaurant/database"
	"restaurant/models"

	"github.com/gin-gonic/gin"
)

// Getables retrieves details of all the tables
func GetTables(c *gin.Context) {
	//reterive the tableinformation
	var tables []models.TablesModel
	database.DB.Find(&tables)
	c.JSON(http.StatusOK, gin.H{
		"Status":  "Success",
		"message": "Table details fetched successfully",
		"data":    tables,
	})
}

func GetTable(c *gin.Context) {
	tableId := c.Param("id")
	var tables models.TablesModel

	if err := database.DB.First(&tables, "table_id =? ", tableId).Error; err != nil {
		c.JSON(400, gin.H{
			"status":  "Failed",
			"message": "Table not found",
			"data":    err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"status":  "Success",
		"message": "Table found succeesfully",
		"data":    tables,
	})

}

// Reserve  a new table with user
func ReserveTable(c *gin.Context) {
	var table models.TablesModel
	if err := c.ShouldBindJSON(&table); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingTable models.TablesModel
	if err := database.DB.First(&existingTable, table.TableID).Error; err != nil {
		c.JSON(404, gin.H{"error": "table not found"})

	}

	//Check if table is available for reservation
	if !IsTableAvailable(uint(table.TableID)){
		c.JSON(http.StatusConflict, gin.H{"error":"The table already booked"})
		return
	}

	//Create Reservation for table
	if err := database.DB.Create(&table).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "table reserved successfully",
		"data":    table,
	})
}

func IsTableAvailable(tableID uint)bool{
	var reservation models.TablesModel

	if err := database.DB.Where("table_id = ?",tableID).First(&reservation).Error; err != nil{
		return true
	}
	return false
}

// Delete the table information with admin authentication
func DeleteTable(c *gin.Context) {
	tableID := c.Param("id")
	var table models.TablesModel
	if err := database.DB.First(&table, tableID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "Failed",
			"message": "Table Not Found",
			"data":    err.Error(),
		})
		return
	}

	database.DB.Delete(&table)
	c.JSON(200, gin.H{
		"status":  "Success",
		"message": "Table deleted successfully",
		"data":    tableID,
	})
}
