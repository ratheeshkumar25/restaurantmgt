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
}

// RazorPay Model
type RazorPay struct {
	UserID          uint    `JSON:"userid"`
	RazorPaymentID  string  `JSON:"razorpaymentid" gorm:"primaryKey"`
	RazorPayOrderID string  `JSON:"razorpayorderid"`
	Signature       string  `JSON:"signature"`
	AmountPaid      float64 `JSON:"amountpaid"`
}
