package models

import "gorm.io/gorm"

type NotesModel struct {
	gorm.Model
	Notes   string `json:"notes"`
	Title   string `json:"title"`
	TableID uint
}
