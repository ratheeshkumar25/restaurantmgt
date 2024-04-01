package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"restaurant/database"
	"restaurant/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Getstaff details
func GetStaff(c *gin.Context) {
	var staff []models.StaffModel
	database.DB.Find(&staff)
	c.JSON(200, gin.H{
		"Status":  "Success",
		"Message": "Staff Details are fetched successfully",
		"data":    staff,
	})
}

// GetStaffByID retrieve a staff member by ID
func GetStaffByIDs(c *gin.Context) {
	var staff models.StaffModel
	staffID := c.Param("id")

	if err := database.DB.First(&staff, staffID).Error; err != nil {
		c.JSON(404, gin.H{"error": "Staff not found"})
		return
	}

	c.JSON(200, gin.H{
		"status":  "Success",
		"message": "Staff details fetched successfully",
		"data":    staff,
	})

}

// Add new staff member to restaurant
func AddStaff(c *gin.Context) {
	var staff models.StaffModel
	if err := c.BindJSON(&staff); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Create(&staff).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{
		"Status":  "Success",
		"Message": "Successfully added staff details",
		"data":    staff,
	})
}

// Update a staff member
func UpdateStaff(c *gin.Context) {
	var staff models.StaffModel
	staffID := c.Param("id")

	if err := database.DB.First(&staff, staffID).Error; err != nil {
		c.JSON(404, gin.H{"error": "Staff not found"})
		return
	}
	if err := c.BindJSON(&staff); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Save(&staff).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{
		"status":  "Success",
		"message": "Staff details updated successfully",
		"data":    staff,
	})

}

// Remove a staff member
func RemoveStaff(c *gin.Context) {
	staffID := c.Param("id")
	var staff models.StaffModel

	if err := database.DB.First(&staff, staffID).Error; err != nil {
		c.JSON(404, gin.H{
			"status":  "Failed",
			"message": "Staff id not found",
			"data":    err.Error(),
		})
		return
	}
	database.DB.Delete(&staff)
	c.JSON(200, gin.H{
		"status":  "success",
		"message": "staff details removed successfully",
	})
}

func StaffAssignTable(c *gin.Context) {
	var request struct {
		StaffID uint `json:"staffID"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(request)
	//Fetch the staff details

	var staff models.StaffModel
	if err := database.DB.First(&staff, request.StaffID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Staff not found"})
		return
	}
	//check staff role and assign table accordingly
	var tableID uint
	switch staff.Role {
	case "Waiter":
		tableID = 1
	case "Waiter1":
		tableID = 2
	case "Waiter2":
		tableID = 3
	case "Waiter3":
		tableID = 4
	case "Waiter4":
		tableID = 5
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid staff role"})
		return
	}
	//Update the staff TableID in the database
	staff.TableID = tableID
	if err := database.DB.Save(&staff).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign table to staff"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "Success",
		"message": "Table assigned to staff successfully",
		"data": gin.H{
			"staff":    staff,
			"table_id": tableID,
		},
	})

}

func fetchStaffIDByTableID(tableID int) (int, error) {
	var staff models.StaffModel
	if err := database.DB.Where("table_id = ?", tableID).First(&staff).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, fmt.Errorf("no staff found for table with ID %d", tableID)
		}
		return 0, fmt.Errorf("failed to fetch staff: %v", err)
	}

	return int(staff.ID), nil
}
