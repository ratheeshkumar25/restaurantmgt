package models

import (
	"time"

	"gorm.io/gorm"
)

//table entity creation

type TablesModel struct {
	gorm.Model
	Capacity     int  `json:"capacity" validate:"required"`
	Availability bool `json:"availability" validate:"required" `
}

type ReservationModels struct {
	gorm.Model
	Date          time.Time `json:"date" gorm:"column:date"`
	TableID       int       `json:"tableID"`
	NumberOfGuest int       `json:"numberofGuest"`
	StartTime     time.Time `json:"startTime"`
	EndTime       time.Time `json:"endTime"`
	UserID        uint      `json:"userID"`
	StaffID 	  uint 		`gorm:"not null"`
}

// BeforeCreate hook to set the Date field to the current date before creating a new record.
func (reservation *ReservationModels) BeforeCreate(tx *gorm.DB) (err error) {
	reservation.Date = time.Now()
	return nil
}
