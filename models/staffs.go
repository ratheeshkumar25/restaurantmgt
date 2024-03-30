package models

import "gorm.io/gorm"

type StaffModel struct {
	gorm.Model
	Staff_name string `json:"staffname"`
	Role       string `json:"staffrole"`
	Salary     int    `json:"salary"`
	TableID    uint
}
