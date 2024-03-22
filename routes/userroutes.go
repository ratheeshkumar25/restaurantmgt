package routes

import (
	"restaurant/controllers"
	"restaurant/middleware"
	//"restaurant/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes() *gin.Engine {
	//create a new engine with default settings
	r := gin.Default()

	//Define user routers
	r.GET("/users", controllers.GetLogin)
	r.POST("users/login", controllers.PostLoginHander)
	r.POST("users/login/verify", controllers.SignupVerify)
	r.POST("/logout", controllers.UserLogout)

	//Define the Admin Routes
	r.POST("/adminlogin/auth", controllers.AdminLogin)
	r.POST("/admin/login", controllers.AdminLogin)
	r.POST("admin/logout", controllers.AdminLogout)

	//admin middleware authentication to add edit and delete
	admin := r.Group("/admin")
	admin.Use(middleware.AdminAuthMiddleware())
	admin.POST("/menu/add", controllers.CreateMenu)
	admin.PUT("menu/:id", controllers.UpdateMenu)
	admin.DELETE("menu/:id", controllers.DeleteMenu)
	admin.GET("/table", controllers.GetTables)
	admin.DELETE("table/:id", controllers.DeleteTable)
	

	//Users middleware authentication view menulist , specified menu
	users := r.Group("/users")
	users.Use(middleware.UserauthMiddleware())
	users.GET("/menulist",controllers.GetMenuList)
	users.POST("table/add", controllers.CreateTable)

	
	

	return r
}

