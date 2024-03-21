package database

import (
	"os"
	"restaurant/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func DBconnect(){
	dsn := os.Getenv("DSN")

	db,err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil{
		panic("failed to connect to database")
	}
	DB = db

	DB.AutoMigrate(&models.UsersModel{},&models.AdminModel{},&models.InvoicesModel{},
	&models.MenuModel{},&models.NotesModel{},&models.NotificationModel{},&models.StaffModel{},&models.TablesModel{},&models.VerifyOTP{},)

}
