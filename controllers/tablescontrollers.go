package controllers

import (
	"net/http"
	"restaurant/database"
	"restaurant/models"

	"github.com/gin-gonic/gin"
)

//Getables retrieves details of all the tables 
func GetTables(c *gin.Context) {
	//reterive the tableinformation 
	var tables []models.TablesModel
	database.DB.Find(&tables)
	c.JSON(http.StatusOK,gin.H{
		"Status":"Success",
		"message":"Table details fetched successfully",
           "data": tables,
		})
}

//Create a new table with user 
func CreateTable(c *gin.Context){
	var table models.TablesModel
	if err := c.ShouldBindJSON(&table); err !=nil{
		c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
	    return
	}
	if err := database.DB.Create(&table).Error; err != nil{
		c.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
	}
c.JSON(http.StatusCreated,table)
}

//Delete the table information with admin authentication 
func DeleteTable(c*gin.Context){
	var table models.TablesModel
	tableID := c.Param("id")
	if err := database.DB.First(&table, tableID).Error; err !=nil{
	c.JSON(http.StatusNotFound,gin.H{
		"status":"Failed",
		"message":"Table Not Found",
		"data": err.Error(),
	})
	return
	}

	database.DB.Delete(&table)
	c.JSON(200,gin.H{
		"status":"Failed",
		"message":"Table Not Found",
		"data": tableID,
	})
}

