package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"restaurant/database"
	"restaurant/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

//Define constants for the payment status

const(
	PaymentPending = "Pending"
	PaymentComplete = "Completed"
)

// view invoice the generated invocie
func GetInvoice(c *gin.Context) {
	var invoice []models.InvoicesModel

	if err := database.DB.Find(&invoice).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "error occured while receiving the invoice "})
		return
	}
	c.JSON(http.StatusOK, invoice)
}

// place the order and generating invoice
func PlaceOrder(c *gin.Context) {
	var invoice models.InvoicesModel

	//bind Json data to the invoice struct 
	if err := c.ShouldBindJSON(&invoice); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//Perform the filed validation 
	if err := validateInvoice(invoice); err != nil{
		c.JSON(400,gin.H{"error":err.Error()})
	}
	//Fetch the menu details from database based on provide MenuID
	menu,err := database.GetMenuByID(invoice.MenuID)
		if err != nil {
			c.JSON(500,gin.H{"error":"Failed to fetch menu detalis"})
			return
		}
	// Set the unit price in the invoice based on menu details 
	invoice.UnitPrice = menu.Price
	//calculate the total amount
	invoice.TotalAmount = float64(invoice.Quantity) * invoice.UnitPrice
	//Set payment due date
	invoice.PaymentDueDate = time.Now().AddDate(0, 0, 7)
	//Set the payment status to paymentPending 
	invoice.PaymentStatus = PaymentPending

	//Automatically fetch staffID based on selected table
	staffID, err := fetchStaffIDByTableID(invoice.TableID)
	if err != nil{
		c.JSON(http.StatusInternalServerError,gin.H{"error":"Failed to fetch staffID"})
		return
	}
	invoice.StaffID = staffID


	//check the order already exists in db
	existingOrder, err := database.GetOrderByID(uint(invoice.OrderID))
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(400, gin.H{"error": "Failed to check the order"})
		return
	}

	//if the orderexist notify the user

	if existingOrder != nil {
		c.JSON(400, gin.H{"error": "Order already exists"})
		return
	}
	// Generate Invoice
	if err := database.DB.Create(&invoice).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create invoice"})
		return
	}
	c.JSON(http.StatusOK,gin.H{
		"status":  "Success",
		"message": "Placed your order successfully",
		"data":    invoice,
	})
}

//Update the placeorder for users 

func UpdatePlaceOrder(c * gin.Context){
	var updateOrder models.InvoicesModel

	if err := c.ShouldBindJSON(&updateOrder);err != nil{
		c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
		return
	}
	updateInvoice := c.Param("id")
	fmt.Print(updateInvoice)
	var existingPlaceOrder models.InvoicesModel

	if err := database.DB.First(&existingPlaceOrder,updateInvoice).Error; err !=nil{
		c.JSON(http.StatusNotFound,gin.H{"error":"Order not found"})
		return
	}

		//Fetch the menu details from database based on provide MenuID
		menu,err := database.GetMenuByID(existingPlaceOrder.MenuID)
		if err != nil {
			c.JSON(500,gin.H{"error":"Failed to fetch menu detalis"})
			return
		}
	// update the filed of existing Placeorder 
    updateOrder.UnitPrice = menu.Price
	updateOrder.TotalAmount = float64(updateOrder.Quantity)*updateOrder.UnitPrice
	updateOrder.PaymentStatus = PaymentPending
	
	existingPlaceOrder.InvoiceID = updateOrder.InvoiceID
	existingPlaceOrder.OrderID = updateOrder.OrderID
	existingPlaceOrder.TableID = updateOrder.TableID
	existingPlaceOrder.Quantity = updateOrder.Quantity
	existingPlaceOrder.MenuID = updateOrder.MenuID
	existingPlaceOrder.UnitPrice = updateOrder.UnitPrice
	existingPlaceOrder.TotalAmount = updateOrder.TotalAmount
	existingPlaceOrder.PaymentDueDate = time.Now().AddDate(0, 0, 7)
	existingPlaceOrder.PaymentStatus = updateOrder.PaymentStatus
	

	//Save the update palceorder

	if err := database.DB.Save(&existingPlaceOrder).Error; err != nil{
		c.JSON(http.StatusInternalServerError,gin.H{"error":"failed to update the order"})
		return
	}
	c.JSON(http.StatusOK,gin.H{
		"status":  "Success",
		"message": "PlaceOrder Details Updated successfully",
		"data":    updateOrder,
	})

}


//Validate Invoice filed 
func validateInvoice(invoice models.InvoicesModel) error {
	if invoice.TableID <= 0 {
		return errors.New("table ID must be positive")
	}
	if invoice.Quantity <= 0 {
		return errors.New("quantity must be positive")
	}
	if invoice.MenuID <= 0 {
		return errors.New("menu ID must be positive")
	}
	return nil
}
//Payinvoice handles payment for an invoice
func PayInvoice(c *gin.Context){
	invoiceID := c.Param("id")
	id,err := strconv.Atoi(invoiceID)
	if err != nil{
		c.JSON(http.StatusBadRequest,gin.H{"error":"Invalid invoice ID"})
		return
	}
	//Fetch the invoice from database

	var invoice models.InvoicesModel
	if err := database.DB.First(&invoice,id).Error;err != nil{
		if errors.Is(err,gorm.ErrRecordNotFound){
			c.JSON(404,gin.H{"error":"Invoice not found"})
			return
		}
		c.JSON(500,gin.H{"error":"Failed to fetch the invoice"})
	return
	}

	//Check if invoice is alredy paid 
	if invoice.PaymentStatus == PaymentComplete{
		c.JSON(400,gin.H{"error":"Invoice is already paid"})
		return 
	}
	//Simulate payment processing
	time.Sleep(3 *time.Second)

	// Update the payment status to completed
	invoice.PaymentStatus = PaymentComplete
	if err := database.DB.Save(&invoice).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update payment status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment successful", "invoice": invoice})
}

	

