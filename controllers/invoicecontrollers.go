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

const (
	PaymentPending  = "Pending"
	PaymentComplete = "Completed"
	PaymentCancelled ="Cancelled"
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
	if err := validateInvoice(invoice); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	}

	var menudetails struct{
		ID int 
		Category string
		Price float64
		FoodImage string
	    Duration string 
	}
	//Fetch the menu details from database based on provide MenuID
	menu, err := database.GetMenuByID(invoice.MenuID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch menu detalis"})
		return
	}
    menudetails.ID = int(menu.ID)
	menudetails.Category = menu.Category
	menudetails.Price = menu.Price
	menudetails.FoodImage = menu.FoodImage
	menudetails.Duration = menu.Duration

	//fetch the reservation details based on the tableId
	reservation, err := database.GetReservationByID(uint(invoice.TableID))
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch reservation details"})
		return
	}

	fmt.Println("hi",reservation)

	// Automatically fetch user ID based on user details
// user, err := database.GetUsersByID(invoice.UserID)
// fmt.Println(user)
// if err != nil {
//     if errors.Is(err, gorm.ErrRecordNotFound) {
//         c.JSON(400, gin.H{"error": "User not found"})
//     } else {
//         c.JSON(500, gin.H{"error": "Failed to fetch user details"})
//     }
//     return
// }


// if user == nil {
//     c.JSON(400, gin.H{"error": "User not found"})
//     return
// }

// 	invoice.UserID = uint(user.UserID)
// fmt.Println(invoice.UserID)	
	//invoice.Quantity = reservation.NumberOfGuest
	// Set the unit price in the invoice based on menu details
	invoice.UnitPrice = menu.Price
	//calculate the total amount
	invoice.TotalAmount = float64(invoice.Quantity) * invoice.UnitPrice
	//Set payment due date
	invoice.PaymentDueDate = time.Now().AddDate(0, 0, 7)
	//Set the payment status to paymentPending
	invoice.PaymentStatus = PaymentPending
	

	//Automatically fetch staffID based on selected table
	// staffID, err := fetchStaffIDByTableID(int(invoice.TableID))
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch staffID"})
	// 	return
	// }
	// invoice.StaffID = staffID

	//Automatically gert access to place order with login customer
	userIDContext, _ := c.Get("userID")
	//fmt.Println(userIDContext)
	userID := userIDContext.(uint)

	var bookingID models.UsersModel

	if err := database.DB.Where("user_id = ?", userID).First(&bookingID).Error; err == gorm.ErrRecordNotFound {
		c.JSON(200, gin.H{"status": "Success",
			"message": "No booking",
			"data":    nil})
		return
	} else if err != nil {
		c.JSON(404, gin.H{"status": "Failed",
			"message": "Database error",
			"data":    nil})
		return
	}
	invoice.UserID = userID

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
	c.JSON(http.StatusOK, gin.H{
		"status":  "Success",
		"message": "Placed your order successfully",
		// "data":    invoice,
		"data":gin.H{
			"data":invoice,
			"menu":menudetails,
		},
	})
}

// Update the placeorder for users
func UpdatePlaceOrder(c *gin.Context) {
	var updateOrder models.InvoicesModel
	if err := c.ShouldBindJSON(&updateOrder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// extract the invoiceid from the request paramenters
	invoiceID := c.Param("id")
	//Fetch the menu  details from database based on existing order
	var existingPlaceOrder models.InvoicesModel
	if err := database.DB.First(&existingPlaceOrder, invoiceID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	//Fetch the menu details from database based on provide MenuID
	menu, err := database.GetMenuByID(existingPlaceOrder.MenuID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch menu detalis"})
		return
	}
	// update the filed of existing Placeorder
	updateOrder.UnitPrice = menu.Price
	updateOrder.TotalAmount = float64(updateOrder.Quantity) * updateOrder.UnitPrice
	updateOrder.PaymentStatus = PaymentPending
	existingPlaceOrder.PaymentMethod = updateOrder.PaymentMethod
	existingPlaceOrder.InvoiceID = updateOrder.InvoiceID
	existingPlaceOrder.OrderID = updateOrder.OrderID
	existingPlaceOrder.TableID = updateOrder.TableID
	existingPlaceOrder.Quantity = updateOrder.Quantity
	existingPlaceOrder.MenuID = updateOrder.MenuID
	existingPlaceOrder.UnitPrice = updateOrder.UnitPrice
	existingPlaceOrder.TotalAmount = updateOrder.TotalAmount
	existingPlaceOrder.PaymentDueDate = time.Now().AddDate(0, 0, 7)
	existingPlaceOrder.PaymentStatus = updateOrder.PaymentStatus

	//Save the update palceorder details back to the database

	if err := database.DB.Save(&existingPlaceOrder).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update the order"})
		return
	}

	//Respond with the updated order details
	c.JSON(http.StatusOK, gin.H{
		"status":  "Success",
		"message": "PlaceOrder Details Updated successfully",
		"data":    updateOrder,
	})

}
//Cancels the existing order
func CancelOrder(c *gin.Context){
	//Extract the invoice id from request parameter
	invoiceID := c.Param("id")

	//Fetch the invoice 
	var invoice models.InvoicesModel
	if err := database.DB.First(&invoice,invoiceID).Error; err != nil{
		c.JSON(404,gin.H{"error":"Order not found"})
		return
	}

	if invoice.PaymentStatus == PaymentCancelled{
		c.JSON(404,gin.H{"error":"Order is already canceled"})
	}


		invoice.PaymentStatus = PaymentCancelled
		if err := database.DB.Save(&invoice).Error; err !=nil{
			c.JSON(500,gin.H{"error":"Failed to cancel the order"})
			return
		}
	// Respond with updated invoce 
	response := gin.H{
		"userID":invoice.UserID,
		"invoiceID":invoice.InvoiceID,
		"menuID":invoice.MenuID,
		"orderID":invoice.OrderID,
		"quantity":invoice.Quantity,
		"totatAmount":invoice.TotalAmount,
		"paymentMethod":invoice.PaymentMethod,
		"paymentstatus":invoice.PaymentStatus,
	}
	c.JSON(200,gin.H{"message":"Order canceled successfully","response":response})

}

// Validate Invoice filed
func validateInvoice(invoice models.InvoicesModel) error {
	if invoice.TableID <= 0 {
		return errors.New("table ID must be positive")
	}
	if invoice.MenuID <= 0 {
		return errors.New("menu ID must be positive")
	}
	return nil
}

// Payinvoice handles payment for an invoice
func PayInvoice(c *gin.Context) {
	invoiceID := c.Param("id")
	id, err := strconv.Atoi(invoiceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid invoice ID"})
		return
	}
	//Fetch the invoice from database

	var invoice models.InvoicesModel
	if err := database.DB.First(&invoice, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(404, gin.H{"error": "Invoice not found"})
			return
		}
		c.JSON(500, gin.H{"error": "Failed to fetch the invoice"})
		return
	}

	//Check if invoice is alredy paid
	if invoice.PaymentStatus == PaymentComplete {
		c.JSON(400, gin.H{"error": "Invoice is already paid"})
		return
	}
	//Simulate payment processing
	//time.Sleep(3 * time.Second)

	// Update the payment status to completed
	invoice.PaymentStatus = PaymentComplete
	if err := database.DB.Save(&invoice).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update payment status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment successful", "invoice": invoice})
}
