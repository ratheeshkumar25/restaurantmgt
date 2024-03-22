package controllers

import (
	"fmt"
	"net/http"
	"restaurant/database"
	"restaurant/models"

	"github.com/gin-gonic/gin"
)
//User can access whole menu list 
func GetMenuList(c *gin.Context) {
	var menus []models.MenuModel
	database.DB.Find(&menus)
	//fmt.Println(menus)
		c.JSON(200, gin.H{
			"status":"Success",
			"message": "Menu details fetched successfully",
			"data":menus,
		})
}
//Access to particular menu for user to check with food id 
func GetMenu(c *gin.Context){
	//Reterive the food_id parameter from the URL 
	foodID := c.Param("id")

	//Declare a slice to hold menuitems 
	var menuItem models.MenuModel

	//Find all items matched with that match the food_id 
	if err := database.DB.First(&menuItem,"food_id = ?",foodID).Error; err != nil{
		c.JSON(404,gin.H{
			"status": "failed",
			"message": "unable found items",
			"data": err.Error(),

		})
		return
	}
		c.JSON(200,gin.H{
			"status":"Success",
			"message":"Menu details fetched successfully",
			"data":menuItem,
		})
}

//Create menulist for admin with authentication 
func CreateMenu(c *gin.Context){
	var menu models.MenuModel
	if err := c.ShouldBindJSON(&menu); err !=nil{
		c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
	    return
	}
	if err := database.DB.Create(&menu).Error; err != nil{
		c.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
	}
	fmt.Println("Menu added",menu)
c.JSON(201,gin.H{"message":"Item added successfully"})
}

//Updae the menu for admin with authentication 
func UpdateMenu(c *gin.Context){
	var menu models.MenuModel

	if err := c.ShouldBindJSON(&menu); err !=nil{
		c.JSON(400,gin.H{"error":err.Error()})
		return
	}
	menuID := c.Param("food_id")
	fmt.Println(menuID)
	var existingMenu models.MenuModel

	if err :=database.DB.First(&existingMenu,menuID).Error; err !=nil{
		c.JSON(http.StatusNotFound,gin.H{"error":"Menu not Found"})
		return
	}

	//update the fiels of existing menulist 
	existingMenu.Food_id = menu.Food_id
	existingMenu.Category = menu.Category
	existingMenu.Name = menu.Name
	existingMenu.Price = menu.Price
	existingMenu.Food_image = menu.Food_image
	existingMenu.Duration = menu.Duration
	existingMenu.TableID = menu.TableID

	//save the updated menu item to the database 

	if err := database.DB.Save(&existingMenu).Error; err !=nil{
		c.JSON(http.StatusInternalServerError,gin.H{"error":"failed to update menu"})
		return
	}
	c.JSON(http.StatusOK,gin.H{
		"status":"Success",
		"message":"Menu Details Updated successfully",
			"data":menu,
		})
}

//Delete the menu for admin with authentication 
func DeleteMenu(c *gin.Context){
	id := c.Param("id")
	var menu models.MenuModel

	if err := database.DB.First(&menu,id).Error; err != nil{
		c.JSON(http.StatusNotFound,gin.H{
			"status": "Failed",
			"message":"User Not Found",
			"data":err.Error(),
		})
		return
	}
    database.DB.Delete(&menu)
	c.JSON(http.StatusOK,gin.H{
		"status":"Failed",
		"message":"User Deleted Successfully",
		"data": id,
	})

}


