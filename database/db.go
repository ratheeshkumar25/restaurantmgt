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

	DB.AutoMigrate(
		&models.UsersModel{}, 
		&models.AdminModel{}, 
		&models.InvoicesModel{},
		&models.MenuModel{},
		&models.ReviewModel{},
		&models.NotificationModel{}, 
		&models.StaffModel{},
	    &models.TablesModel{},
		&models.RazorPay{},
		&models.ReservationModels{},
	)

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

func GetUsersByID(userID uint) (*models.UsersModel, error) {
    // Return nil if userID is 0
    if userID == 0 {
        return nil, nil
    }

    var user models.UsersModel
    if err := DB.Where("user_id = ?", userID).First(&user).Error; err != nil {
        return nil, err
    }
    return &user, nil
}



func GetReservationByID(tableID uint)(*models.TablesModel,error){
	var reservation models.TablesModel
	if err := DB.Where("table_id = ?",tableID).First(&reservation).Error; err != nil{
		return nil,err
	}
	return &reservation, nil
}