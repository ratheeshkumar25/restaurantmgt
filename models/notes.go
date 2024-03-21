package models

type NotesModel struct {
	Notes_id int    `gorm:"primaryKey"`
	Notes    string `json:"notes"`
	Title    string `json:"title"`
	// Table_Id TablesModel `gorm:"foreignKey:table_id"`
	TableID  uint 
}
