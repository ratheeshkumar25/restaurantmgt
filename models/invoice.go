package models

import "time"

//Invoice Model
type InvoicesModel struct {
	InvoiceID      int       `gorm:"primaryKey;autoIncrement"`
	OrderID        int       `gorm:"autoIncrement"`
	TableID        int       `json:"tableID"`
	StaffID        int       `json:"staffID"`
	Quantity       int       `json:"quantity"`
	UnitPrice      float64   `json:"unitPrice"`
	TotalAmount    float64   `json:"totalAmount"`
	PaymentMethod  string    `json:"paymentMethod"`
	PaymentDueDate time.Time `json:"paymentDueDate"`
	PaymentStatus  string    `json:"paymentStatus"`
	MenuID         uint
	UserID         uint
}

// RazorPay Model
type RazorPay struct {
	InvoiceID       uint    `JSON:"userID"`
	RazorPaymentID  string  `JSON:"razorpaymentID" gorm:"primaryKey;autoIncrement"`
	RazorPayOrderID string  `JSON:"razorpayorderID"`
	Signature       string  `JSON:"signature"`
	AmountPaid      float64 `JSON:"amountpaid"`
}
