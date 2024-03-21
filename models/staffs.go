package models

import "gorm.io/gorm"

type StaffModel struct {
	gorm.Model
	Staff_id   int         `gorm:"primaryKey"`
	Staff_name string      `json:"staffname"`
	Role       string      `json:"staffrole"`
	Salary     int         `json:"salary"`
	// TableID	   TablesModel `gorm:"foreignKey:table_id"`
	TableID    uint 
}

type OrderModel struct{
	gorm.Model
	OrderID int 	`gorm:"primaryKey"`
	TableID uint        
	StaffID int 
	Description string `json:"description"`
}
