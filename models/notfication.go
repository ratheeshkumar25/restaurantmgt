package models

import "time"

type NotificationModel struct {
	Notification_id int 	  `gorm:"primarykey"`
	Invoice_id      uint
	Ordertaken_time time.Time `json:"ordertakentime"`
	End_time        time.Time `json:"endtime"`
	Status          string    `json:"status"`
}
