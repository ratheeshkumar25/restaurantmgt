package controllers

import (
	"net/http"
	"restaurant/database"
	"restaurant/models"

	"github.com/gin-gonic/gin"
)


func GetTables(c *gin.Context) {
	var tables []models.TablesModel
	database.DB.Find(&tables)
	c.JSON(http.StatusOK,gin.H{
		"Status":"Success",
		"message":"Table details fetched successfully",
           "data": tables,
		})
}
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

func DeleteTable(c*gin.Context){
	var table models.TablesModel
	tableID := c.Param("tableID")
	if err := database.DB.First(&table, tableID).Error; err !=nil{
	c.JSON(http.StatusNotFound,gin.H{"error":"Table Not Found"})
	return
	}
}

