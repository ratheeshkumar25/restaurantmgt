package controllers

import (
	"net/http"
	"restaurant/database"
	"restaurant/models"

	"github.com/gin-gonic/gin"
)

func CreateNote(c *gin.Context) {
	var note models.NotesModel

	if err := c.ShouldBindJSON(&note); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Create(&note).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create note"})
		return
	}

	c.JSON(http.StatusCreated, note)
}