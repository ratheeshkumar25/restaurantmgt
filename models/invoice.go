package models

import "time"

type InvoicesModel struct {
	Invoice_id       int       `gorm:"primaryKey"`
	Order_id         int       `json:"orderid"`
	Quantity         int       `json:"quantity"`
	Unit_price       float64   `json:"unit_price"`
	Total_amount     float64   `json:"total_amount"`
	Payment_method   string    `json:"payment_method"`
	Payment_due_date time.Time `json:"payment_due_date"`
	Food_id    		 uint
}