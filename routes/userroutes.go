package routes

import (
	"restaurant/controllers"
	"restaurant/middleware"

	//"restaurant/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes() *gin.Engine {
	//creates a new Gin engine instance with default configurations
	r := gin.Default()

	//Define user routers
	r.GET("/users", controllers.GetLogin)
	r.POST("/users/login", controllers.PostLoginHander)
	r.POST("/users/login/verify", controllers.SignupVerify)
	r.POST("/logout", controllers.UserLogout)

	//Define the Admin Routes
	r.POST("/admin/login", controllers.AdminLogin)
	r.POST("/admin/logout", controllers.AdminLogout)

	//admin middleware authentication to add edit and delete
	admin := r.Group("/admin")
	admin.Use(middleware.AdminAuthMiddleware())
	{
		admin.GET("/menuList", controllers.GetMenuList)
		admin.POST("/menu/add", controllers.CreateMenu)
		admin.PUT("/menu/:id", controllers.UpdateMenu)
		admin.DELETE("menu/:id", controllers.DeleteMenu)
		//table control
		admin.GET("/table", controllers.GetTables)
		admin.GET("/table/:id", controllers.GetTable)
		admin.DELETE("/table/:id", controllers.DeleteTable)
		//staff control
		admin.GET("/staff", controllers.GetStaff)
		admin.GET("/staff/:id", controllers.GetStaffByIDs)
		admin.POST("/staff/add", controllers.AddStaff)
		admin.PUT("/staff/:id", controllers.UpdateStaff)
		admin.POST("/staff/:id", controllers.StaffAssignTable)
		admin.DELETE("/staff/:id", controllers.RemoveStaff)
		//order and invoice controller
		admin.GET("invoice", controllers.GetInvoice)
	}

	//Users middleware authentication view menulist , specified menu
	users := r.Group("/users")
	users.Use(middleware.UserauthMiddleware())
	{
		users.GET("menu/:id", controllers.GetMenu)
		users.GET("/menulist", controllers.GetMenuList)
		users.GET("/table", controllers.GetTables)
		users.POST("/table/add", controllers.ReserveTable)
		users.POST("/placeorder/invoice", controllers.PlaceOrder)
		users.POST("/payinvoice/:id", controllers.PayInvoice)
		users.PUT("/updateorder/:id", controllers.UpdatePlaceOrder)
		//users.GET("invoice",controllers.GetInvoice)
	}

	return r
}
