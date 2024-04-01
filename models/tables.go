package models

import "gorm.io/gorm"

//table entity creation

type TablesModel struct {
	gorm.Model
	TableID         int  `json:"tableID"`
	Number_of_guest int  `json:"numberofGuest"`
	UserID          uint `json:"userID"`
}
