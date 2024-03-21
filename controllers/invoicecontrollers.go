package controllers

import (
	"net/http"
	"restaurant/database"
	"restaurant/models"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateInvoice(c *gin.Context){
	var invoice models.InvoicesModel

	//data from the request body to inv struct 
	if err := c.ShouldBindJSON(&invoice); err !=nil{
		c.JSON(http.StatusBadRequest,gin.H{
			"error": err.Error()})
			return
	}

	//filed validation 
	if invoice.Order_id <= 0{
		c.JSON(http.StatusBadRequest,gin.H{
			"error":"Order ID must be positive"})
			return
	}

	if invoice.Quantity <= 0{
		c.JSON(http.StatusBadRequest,gin.H{
			"error":"Quantity must be positive"})
			return
	}

	if invoice.Unit_price <= 0{
		c.JSON(http.StatusBadRequest,gin.H{
			"error":"Unit ID must be positive"})
			return
	}

	//calculate the total amount 

	invoice.Total_amount = float64(invoice.Quantity) * invoice.Unit_price
	
	invoice.Payment_due_date = time.Now().AddDate(0,0,7)

	if err := database.DB.Create(&invoice).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create invoice"})
		return
	}

	c.JSON(http.StatusCreated, invoice)
}

func GetInvoice(c *gin.Context){
	var invoice []models.InvoicesModel

	if err := database.DB.Find(&invoice).Error; err !=nil{
		c.JSON(http.StatusInternalServerError,gin.H{
			"error":"error occured while receiving the invoice "})
			return
	}
	c.JSON(http.StatusOK,invoice)
}