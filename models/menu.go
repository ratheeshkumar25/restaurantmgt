package models

import "time"

type MenuModel struct {
	Food_id   int       `gorm:"primaryKey"`
	Category   string    `json:"category"`
	Name       string    `json:"name"`
	Price      float64   `json:"price"`
	Food_image string    `json:"foodimage"`
	Duration   time.Time `json:"duration"`
	TableID    uint 
}
