package main

import (
	"tess/controllers"
	"tess/helper"
	"tess/initializers"
	"tess/middleware"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
	initializers.SyncDatabase()
}

func main() {
	r := gin.Default()
	r.Use(gin.Logger())

	//user routes
	r.POST("/signup", controllers.UserSignup)
	r.POST("/login", controllers.UserLogin)
	r.GET("/validate", middleware.RequireAuth, controllers.Validate)
	r.GET("/logout", controllers.UserLogout)

	//admin routes
	r.POST("/Admin/login", controllers.Loginadmin)
	r.GET("/Admin/Validate", middleware.AdminAuth, controllers.AdminValidate)
	r.POST("/Admin/logout", controllers.AdminLogout)

	//Create Admin DELETE AFTER SAVE
	r.POST("/adminCreate", helper.AdminCreate)

	//CRUD operation
	r.GET("admin/findall", middleware.AdminAuth, controllers.FindAll)
	r.GET("admin/finduser", middleware.AdminAuth, controllers.FindUser)
	r.PATCH("admin/updateuser", middleware.AdminAuth, controllers.EditUser)
	r.DELETE("admin/deleteuser", middleware.AdminAuth, controllers.DeleteUser)

	r.Run() // listen and serve on 0.0.0.0:8080 start the server on localhost:8080

}
