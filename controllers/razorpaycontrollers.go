package controllers

import (
	"errors"
	"fmt"
	"os"
	//"os/user"
	"restaurant/database"
	//"restaurant/middleware"
	"restaurant/models"
	"strconv"

	//"time"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/razorpay/razorpay-go"

	//"github.com/twilio/twilio-go/client"
	"gorm.io/gorm"
)


type PageVariable struct{
	OrderID string
}
func MakePayment(c *gin.Context){
 // Extract invoice ID from the query parameter
	invoiceID :=c.Query("id")
	id,err := strconv.Atoi(invoiceID)
	if err != nil{
		c.JSON(400,gin.H{"errror":"Invalid invoice ID"})
	}
	var invoice models.InvoicesModel
	if err :=database.DB.First(&invoice,id).Error;err != nil{
		if errors.Is(err,gorm.ErrRecordNotFound){
			c.JSON(404,gin.H{
				"status":"Failed",
				"messsage":"Invoice not found",
				"data":err.Error(),
			})
			return
		}
		c.JSON(500,gin.H{"error":"Failed to fetch the invoice"})
		return 
	}

	  // Check if the invoice is already paid
	  if invoice.PaymentStatus == "Completed" {
        c.JSON(400, gin.H{"error": "Invoice is already paid"})
        return
    }
	// paylink := "pay"
	// if err := database.Client.Set(ctx,"invoice"+invoiceID,paylink,15*time.Minute).Err(); err != nil{
	// 	c.JSON(500,gin.H{
	// 		"status":"Failed",
	// 		"message":"Error setting data in redis server",
	// 		"data":err.Error(),})
	// 		return
	// }


	razorpayment := &models.RazorPay{
		InvoiceID : uint(invoice.InvoiceID),
		AmountPaid: invoice.TotalAmount,
	}
	razorpayment.RazorPaymentID = generateUniqueID()
fmt.Println(razorpayment)
	if err := database.DB.Create(&razorpayment).Error; err != nil{
		c.JSON(500,gin.H{"error":"Failed to create Razorpay Payement"})
		return
	}

	amountInPaisa := invoice.TotalAmount*100

	razorpayClient := razorpay.NewClient(os.Getenv("RAZORPAY_KEY_ID"),os.Getenv("RAZORPAY_SECRET"))

	data := map[string]interface{}{
		"amount":amountInPaisa,
		"currency":"INR",
		"receipt":"some_receipt_id",

	}
	body, err := razorpayClient.Order.Create(data,nil)

	if err != nil{
		// fmt.Println("Problem getting repository information:%v\n",err)
		// os.Exit(1)
		c.JSON(500,gin.H{"error":"Failed to create razorpay order"})
	}
	value := body["id"]
	str := value.(string)

	homepageVariables := PageVariable{
		OrderID: str,
	}
	c.HTML(200,"app.html",gin.H{
		"invoiceID":id,
		"totalPrice":amountInPaisa/100,
		"total": amountInPaisa,
		"orderID":homepageVariables.OrderID,
		
	})
	
}

// SuccessPage renders the success page.
func SuccessPage(c *gin.Context) {
	pID := c.Query("id")
	fmt.Println(pID)
	fmt.Println("Fully successful")

	var invoice models.InvoicesModel
	if err := database.DB.First(&invoice,pID).Error;err !=nil{
		c.JSON(500,gin.H{"error":"Failed to fetch the invoice"})
		return
	}

	if invoice.PaymentStatus == "Pending"{
		if err := database.DB.Model(&invoice).Update("payment_status","Completed").Error;err != nil{
			c.JSON(500,gin.H{"error": "Failed to update invoice payment status"})
			return
		}
	}

	razorPayment := models.RazorPay{
		InvoiceID: uint(invoice.InvoiceID),
		RazorPaymentID: generateUniqueID(),
		AmountPaid: invoice.TotalAmount,
	}

	if err := database.DB.Create(&razorPayment).Error; err !=nil{
		c.JSON(500,gin.H{"error":"Failed to create Razorpay payment"})
	}
	//Render the successpage 
	c.HTML(200, "success.html", gin.H{
		"paymentID": pID,
		"amountPaid":invoice.TotalAmount,
		"invoiceID":invoice.InvoiceID,
	})

}

// generateUniqueID generates a unique ID using UUID.
func generateUniqueID() string {
    // Generate a Version 4 (random) UUID
    id := uuid.New()
    return id.String()
}