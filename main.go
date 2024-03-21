package main

import (
	"restaurant/database"
	"restaurant/helper"
	"restaurant/routes"
	//"github.com/gin-gonic/gin"
)

func Init(){
	helper.LoadEnv()
	database.DBconnect()
	database.InitRedis()
}

func main() {

	Init()
	r := routes.UserRoutes()
	

	//Run the engine the port 3000
	if err := r.Run(":3000"); err !=nil{
		panic(err) //Handle error if unable to start server
	}

}