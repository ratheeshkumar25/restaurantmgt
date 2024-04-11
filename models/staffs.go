package models

import "gorm.io/gorm"

type StaffModel struct {
	gorm.Model
	StaffName string `json:"staffname"`
	Role       string `json:"staffrole"`
	Salary     int    `json:"salary"`
	Blocked    bool   `json:"blocked"`
	//TableID    uint
}
