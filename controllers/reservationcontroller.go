package controllers

import (
	"fmt"
	"restaurant/database"
	"restaurant/models"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SearchAvailableTables(c *gin.Context) {
	var request struct {
		StartTime     time.Time `json:"startTime"`
		EndTime       time.Time `json:"endTime"`
		NumberOfGuest int       `json:"numberOfGuest"`
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var availableTables []models.TablesModel
	err := database.DB.Where("capacity >= ? AND availability = true", request.NumberOfGuest).Find(&availableTables).Error
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Filter out the tables that are already reserved for the requested time slot
	filteredTables := make([]models.TablesModel, 0)
	for _, table := range availableTables {
		var reservation models.ReservationModels
		err = database.DB.Where("table_id = ? AND ((start_time <= ? AND end_time >= ?) OR (start_time <= ? AND end_time >= ?) OR (start_time >= ? AND end_time <= ?))", table.ID, request.StartTime, request.StartTime, request.EndTime, request.EndTime, request.StartTime, request.EndTime).First(&reservation).Error
		if err == gorm.ErrRecordNotFound {
			filteredTables = append(filteredTables, table)
		}
	}

	c.JSON(200, gin.H{
		"message": "Available tables fetched successfully",
		"data":    filteredTables,
	})
}

func CreateReservartion(c *gin.Context) {
	var reservation models.ReservationModels
	if err := c.BindJSON(&reservation); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	//check for available tables and staff
	availableTable, availableStaff, err := checkAvailability(reservation.StartTime, reservation.EndTime, reservation.NumberOfGuest)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	//Assign tthe available table and staff to reservation
	reservation.TableID = int(availableTable.ID)
	reservation.StaffID = availableStaff.ID
	//context
	userIDContext, _ := c.Get("userID")
	fmt.Println(userIDContext)
	userID := userIDContext.(uint)

	var bookingID models.UsersModel

	if err := database.DB.Where("user_id = ?", userID).First(&bookingID).Error; err == gorm.ErrRecordNotFound {
		c.JSON(200, gin.H{"status": "Success",
			"message": "No booking",
			"data":    nil})
		return
	} else if err != nil {
		c.JSON(404, gin.H{"status": "Failed",
			"message": "Database error",
			"data":    nil})
		return
	}
	reservation.UserID = userID
	if err := database.DB.Create(&reservation).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, gin.H{
		"message":     "Reservation created successfully",
		"reservation": reservation.ID,
		"customerID":reservation.UserID,
		"table":       reservation.TableID,
		"startTime":   reservation.StartTime,
		"endTime":     reservation.EndTime,
		"staffServe":  reservation.StaffID,
	})

}

func UpdateReservation(c *gin.Context) {
	reservationID := c.Param("id")
	var reservation models.ReservationModels

	if err := database.DB.First(&reservation, reservationID).Error; err != nil {
		c.JSON(404, gin.H{"error": "Reservation not found"})
		return
	}

	var updatedReservation models.ReservationModels
	if err := c.BindJSON(&updatedReservation); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Check for available tables and staff members for the new time slot
	availableTable, availableStaff, err := checkAvailability(updatedReservation.StartTime, updatedReservation.EndTime, reservation.NumberOfGuest)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Update the reservation with the new table and staff assignments
	reservation.TableID = int(availableTable.ID)
	reservation.StaffID = availableStaff.ID
	reservation.StartTime = updatedReservation.StartTime
	reservation.EndTime = updatedReservation.EndTime

	if err := database.DB.Save(&reservation).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "Reservation updated successfully",
		"data":    reservation,
	})
}

func CancelReservation(c *gin.Context) {
	reservationID := c.Param("id")
	var reservation models.ReservationModels

	if err := database.DB.First(&reservation, reservationID).Error; err != nil {
		c.JSON(404, gin.H{"error": "Reservation not found"})
		return
	}

	if err := database.DB.Delete(&reservation).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "Reservation canceled successfully",
	})
}

// func checkAvailability(startTime, endTime time.Time, numGuests int) (*models.TablesModel, *models.StaffModel, error) {
//     var availableTable models.TablesModel
//     var availableStaff models.StaffModel

//     // Check for available tables with sufficient capacity
//     err := database.DB.Where("capacity >= ? AND availability = true", numGuests).Find(&availableTable).Error
//     if err != nil {
//         return nil, nil, fmt.Errorf("no available tables found")
//     }

   
// 	//Get and check available staff 
// 	var availableStaffMember []models.StaffModel
// 	err = database.DB.Where("blocked = false").Find(&availableStaff).Error
//     if err != nil {
//         return nil, nil, fmt.Errorf("no available staff found")
//     }
//     // Check if the table and staff are not already reserved for the requested duration
//     for _, staff := range availableStaffMember {
//         var reservation models.ReservationModels
//         err = database.DB.Where("table_id = ? AND staff_id = ? AND ((start_time <= ? AND end_time >= ?) OR (start_time <= ? AND end_time >= ?) OR (start_time >= ? AND end_time <= ?))", availableTable.ID, staff.ID, startTime, startTime, endTime, endTime, startTime, endTime).First(&reservation).Error
//         if err != nil && err != gorm.ErrRecordNotFound {
//             return nil, nil, fmt.Errorf("failed to check availability: %v", err)
//         }

//         if err == gorm.ErrRecordNotFound {
//             availableStaff = staff
//             break
//         }
//     }

//     if availableStaff.ID == 0 {
//         return nil, nil, fmt.Errorf("selected table and staff are not available for the requested time slot")
//     }

//     return &availableTable, &availableStaff, nil
// }

func checkAvailability(startTime, endTime time.Time, numGuests int) (*models.TablesModel, *models.StaffModel, error) {
    var availableTable models.TablesModel
    var availableStaff models.StaffModel

    // Check for available tables with sufficient capacity
    err := database.DB.Where("capacity >= ? AND availability = true", numGuests).First(&availableTable).Error
    if err != nil {
        return nil, nil, fmt.Errorf("no available tables found")
    }

    // Get all available staff members
    var availableStaffMembers []models.StaffModel
    err = database.DB.Where("blocked = false").Find(&availableStaffMembers).Error
    if err != nil {
        return nil, nil, fmt.Errorf("no available staff found")
    }

    // Check if the table is not already reserved for the requested duration
    var reservations []models.ReservationModels
    err = database.DB.Where("table_id = ? AND ((start_time < ? AND end_time > ?) OR (start_time >= ? AND start_time < ?) OR (end_time > ? AND end_time <= ?))", availableTable.ID, startTime, endTime, startTime, endTime, startTime, endTime).Find(&reservations).Error
    if err != nil {
        return nil, nil, fmt.Errorf("failed to check table availability: %v", err)
    }

    if len(reservations) > 0 {
        return nil, nil, fmt.Errorf("selected table is not available for the requested time slot")
    }

    // Check if a staff member is available for the requested duration
    for _, staff := range availableStaffMembers {
        var reservationsForStaff []models.ReservationModels
        err = database.DB.Where("staff_id = ? AND ((start_time < ? AND end_time > ?) OR (start_time >= ? AND start_time < ?) OR (end_time > ? AND end_time <= ?))", staff.ID, startTime, endTime, startTime, endTime, startTime, endTime).Find(&reservationsForStaff).Error
        if err != nil {
            return nil, nil, fmt.Errorf("failed to check staff availability: %v", err)
        }

        if len(reservationsForStaff) == 0 {
            availableStaff = staff
            break
        }
    }

    if availableStaff.ID == 0 {
        return nil, nil, fmt.Errorf("no staff available for the requested time slot")
    }

    return &availableTable, &availableStaff, nil
}