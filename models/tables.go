package models

import "gorm.io/gorm"

//table entity creation

type TablesModel struct {
	gorm.Model
	TableID         int `gorm:"primaryKey;autoIncrement"`
	Number_of_guest int `json:"numberofguest"`
	UserID          uint
}
