package database

import (
	"os"
	"restaurant/models"
     "errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func DBconnect() {
	dsn := os.Getenv("DSN")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}
	DB = db

	DB.AutoMigrate(&models.UsersModel{}, &models.AdminModel{}, &models.InvoicesModel{},
		&models.MenuModel{}, &models.NotesModel{}, &models.NotificationModel{}, 
		&models.StaffModel{}, &models.TablesModel{}, &models.VerifyOTP{},)

}

func GetOrderByID(orderID uint)(*models.InvoicesModel,error){
	var order models.InvoicesModel
	if err := DB.First(&order,orderID).Error; err != nil{
		if errors.Is(err, gorm.ErrRecordNotFound){
			return nil,nil
		}
		return nil, err
	}
	return &order,nil
}

func GetMenuByID(menuID uint)(*models.MenuModel,error){
	var menu models.MenuModel
	if  err := DB.First(&menu,menuID).Error; err != nil{
		if errors.Is(err,gorm.ErrRecordNotFound){
			return nil,nil
		}
		return nil, err
	}
	return &menu,nil
}
