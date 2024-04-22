package controllers

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"restaurant/database"
	"restaurant/middleware"
	"restaurant/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
	"gorm.io/gorm"
)

//Define constants for the payment status

const (
	PaymentPending   = "Pending"
	PaymentComplete  = "Completed"
	PaymentCancelled = "Cancelled"
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
	// var invoice models.InvoicesModel
	var orderRequest struct {
		Items []struct {
			ItemID   uint `json:"itemID"`
			Quantity int  `json:"quantity"`
		} `json:"items"`
		Email         string `json:"email"`
		PaymentMethod string `json:"paymentMethod"`
	}

	// Bind JSON data to the order request struct
	if err := c.ShouldBindJSON(&orderRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Fetch menu details for each item in the order request
	var orderItems []map[string]interface{}
	var totalAmount float64
	for _, item := range orderRequest.Items {
		menu, err := database.GetMenuByID(item.ItemID)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to fetch menu details"})
			return
		}
		orderItem := map[string]interface{}{
			"menu":     menu.ID,
			"category": menu.Category,
			"price":    menu.Price,
			"quantity": item.Quantity,
		}
		orderItems = append(orderItems, orderItem)
		totalAmount += float64(item.Quantity) * menu.Price
	}

	// Create the invoice
	invoice := models.InvoicesModel{
		Quantity:       len(orderRequest.Items),
		TotalAmount:    totalAmount,
		Email:          orderRequest.Email,
		PaymentMethod:  orderRequest.PaymentMethod,
		PaymentDueDate: time.Now().AddDate(0, 0, 7),
		PaymentStatus:  PaymentPending,
	}

	//Automatically get access to place order with login customer
	userIDContext, _ := c.Get("userID")
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

	//fetch the reservation details based on the tableId
	reservation, err := database.GetReservationByID(userID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch reservation details"})
		return
	}
	invoice.TableID = int(reservation.TableID)
	invoice.StaffID = int(reservation.StaffID)
	invoice.UserID = userID

	firstItem := orderRequest.Items[0]
	menu, err := database.GetMenuByID(firstItem.ItemID)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to fetch menu details"})
	}
	invoice.ItemID = uint(menu.ID)

	// Check if the order already exists in the database
	existingOrder, err := database.GetOrderByID(uint(invoice.OrderID))
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(400, gin.H{"error": "Failed to check the order"})
		return
	}

	// If the order exists, notify the user
	if existingOrder != nil {
		c.JSON(400, gin.H{"error": "Order already exists"})
		return
	}

	// Generate Invoice
	if err := database.DB.Create(&invoice).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create invoice"})
		return
	}

	// Construct the response JSON object including invoice details and selected items
	response := gin.H{
		"message":        "Placed your order successfully",
		"invoiceID":      invoice.InvoiceID,
		"totalAmount":    invoice.TotalAmount,
		"email":          invoice.Email,
		"paymentMethod":  invoice.PaymentMethod,
		"paymentDueDate": invoice.PaymentDueDate,
		"paymentStatus":  invoice.PaymentStatus,
		"userID":         invoice.UserID,
		"staffID":        invoice.StaffID,
		"tableID":        invoice.TableID,
		"orderItems":     orderItems,
	}

	c.JSON(http.StatusOK, response)
}

// Cancels the existing order
func CancelOrder(c *gin.Context) {
	//Extract the invoice id from request parameter
	invoiceID := c.Param("id")

	//Fetch the invoice
	var invoice models.InvoicesModel
	if err := database.DB.First(&invoice, invoiceID).Error; err != nil {
		c.JSON(404, gin.H{"error": "Order not found"})
		return
	}

	if invoice.PaymentStatus == PaymentCancelled {
		c.JSON(404, gin.H{"error": "Order is already canceled"})
	}

	invoice.PaymentStatus = PaymentCancelled
	if err := database.DB.Save(&invoice).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to cancel the order"})
		return
	}
	// Respond with updated invoce
	response := gin.H{
		"userID":        invoice.UserID,
		"invoiceID":     invoice.InvoiceID,
		"itemID":        invoice.ItemID,
		"orderID":       invoice.OrderID,
		"quantity":      invoice.Quantity,
		"totatAmount":   invoice.TotalAmount,
		"paymentMethod": invoice.PaymentMethod,
		"paymentstatus": invoice.PaymentStatus,
	}
	c.JSON(200, gin.H{"message": "Order canceled successfully", "response": response})

}

func UpdatePlaceOrder(c *gin.Context) {
	// Get the order ID from the request parameters
	orderID := c.Param("id")
	fmt.Println("hello", orderID)
	var invoice models.InvoicesModel
	if err := database.DB.First(&invoice, orderID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(404, gin.H{"error": "invoicde not found"})
			return
		}
		c.JSON(500, gin.H{"error": "Failed to fetch invoice"})
		return
	}

	var orderRequest struct {
		Items []struct {
			ItemID   uint `json:"itemID"`
			Quantity int  `json:"quantity"`
		} `json:"items"`
		Email         string `json:"email"`
		PaymentMethod string `json:"paymentMethod"`
	}
	if err := c.ShouldBindJSON(&orderRequest); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if invoice.PaymentStatus == PaymentComplete {
		c.JSON(400, gin.H{"error": "cannot update a completed order"})
		return
	}
	// Fetch menu details for each item in the order request
	var updateorderItems []map[string]interface{}
	var totalAmount float64
	for _, item := range orderRequest.Items {
		menu, err := database.GetMenuByID(item.ItemID)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to fetch menu details"})
			return
		}
		orderItem := map[string]interface{}{
			"item":     menu.ID,
			"category": menu.Category,
			"price":    menu.Price,
			"quantity": item.Quantity,
		}
		updateorderItems = append(updateorderItems, orderItem)
		totalAmount += float64(item.Quantity) * menu.Price
	}

	invoice.Quantity = len(orderRequest.Items)
	invoice.TotalAmount = totalAmount
	invoice.PaymentMethod = orderRequest.PaymentMethod
	invoice.Email = orderRequest.Email

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
	// 	//fetch the reservation details based on the tableId
	reservation, err := database.GetReservationByID(userID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch reservation details"})
		return
	}
	invoice.TableID = int(reservation.TableID)
	invoice.StaffID = int(reservation.StaffID)
	invoice.UserID = userID

	updateItem := orderRequest.Items[0]
	menu, err := database.GetMenuByID(updateItem.ItemID)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to fetcch the menu details"})
	}
	invoice.ItemID = uint(menu.ID)

	if err := database.DB.Save(&invoice).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update invoice"})
		return
	}

	response := gin.H{
		"itemID":      invoice.ItemID,
		"email":       invoice.Email,
		"userID":      invoice.UserID,
		"staffID":     invoice.StaffID,
		"tableID":     invoice.TableID,
		"totalAmount": invoice.TotalAmount,
	}

	c.JSON(200, gin.H{
		"status":           "Success",
		"message":          "Order updated successfully",
		"data":             response,
		"updateOrderItems": updateorderItems,
	})

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

	// err = middleware.SendEmail("Payment successful for invoice"+invoiceID, invoice.Email)
	// if err != nil {
	// 	c.JSON(500, gin.H{"error": "Failed to send email"})
	// 	return
	// }

	//Generate PDF invoice
	pdfBytes, err := GeneratePDFInvoice(invoice)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate PDF invoice"})
		return
	}

		// // Send payment confirmation email with PDF invoice attached
		err = middleware.SendEmail("Payment successful for invoice"+invoiceID, invoice.Email, "invoice.pdf", pdfBytes)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to send email"})
			return
		}
	//  c.JSON(200, gin.H{"messgae": "Email notification send successfully"})

	c.JSON(http.StatusOK, gin.H{"message": "Payment successful", "invoice": invoice})
}

func GetPDFInvoice(c *gin.Context) {
	invoiceID := c.Param("id")
	id, err := strconv.Atoi(invoiceID)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invaild invoice"})
		return
	}
	//Fetch the invoice from database
	var invoce models.InvoicesModel
	if err := database.DB.First(&invoce, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(404, gin.H{"error": "Invoice not found"})
			return
		}
		c.JSON(500, gin.H{"error": "Failed to fetch the invoice"})
	}
	//Generate PDF invoice
	pdfBytes, err := GeneratePDFInvoice(invoce)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate PDF invoice"})
		return
	}

	// Set response headers
	c.Header("Content-Disposition", "attachment; filename=invoice.pdf")
	c.Header("Content-Type", "application/pdf")
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}

func GeneratePDFInvoice(invoice models.InvoicesModel) ([]byte, error) {
    pdf := gofpdf.New("P", "mm", "A6", "")
    pdf.AddPage()
    pdf.SetFont("Arial", "B", 10)
    // pdf.Cell(40, 10, "Invoice")
	// pdf.Ln(20)
    // // Add invoice details to the PDF
    // pdf.Cell(20, 5, fmt.Sprintf("Invoice ID: %d", invoice.InvoiceID))
	// pdf.Cell(20,5,fmt.Sprintf("Quantity:%d",invoice.Quantity))
    // pdf.Cell(30, 10, fmt.Sprintf("Total Amount: %.2f", invoice.TotalAmount))
    // // Add more details as needed
	    // Add "Go Restaurant" as header
		pdf.Cell(40, 10, "Go Restaurant")
		pdf.Ln(10) // Move down for spacing
	
		// Add title "Invoice"
		pdf.SetFont("Arial", "B", 12) // Set font size for the title
		pdf.Cell(40, 10, "Invoice")
		pdf.Ln(10) // Move down for spacing
	
		// Set font size for details
		pdf.SetFont("Arial", "", 10)
	
		// Add invoice details to the PDF
		pdf.Cell(20, 5, fmt.Sprintf("Invoice ID: %d", invoice.InvoiceID))
		pdf.Ln(5)
		pdf.Cell(20, 5, fmt.Sprintf("Quantity: %d", invoice.Quantity))
		pdf.Ln(10) // Move down for spacing
		pdf.Cell(30, 10, fmt.Sprintf("Total Amount: %.2f", invoice.TotalAmount))
	
		// Add line
		pdf.Line(10, pdf.GetY(), 90, pdf.GetY())
	
		// Add "Thank you for choosing. Welcome back again!"
		pdf.Ln(10) // Move down for spacing
		pdf.Cell(40, 10, "Thank you for choosing. Welcome back again!")
	var buf bytes.Buffer
    err := pdf.Output(&buf)
    if err != nil {
        return nil, err
    }

    return buf.Bytes(), nil
}

