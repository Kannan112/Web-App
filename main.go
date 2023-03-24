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
	r := gin.Default() //gin engin with default middleware

	r.POST("/adminCreate", helper.AdminCreate)

	r.POST("/signup", controllers.UserSignup)

	r.POST("/login", controllers.UserLogin)

	r.GET("/validate", middleware.RequireAuth, controllers.Validate)

	r.PATCH("/updateuser", middleware.RequireAuth, controllers.EditUser)

	r.POST("/logout", controllers.UserLogout)

	r.POST("/Adminlogin", controllers.Loginadmin)

	r.GET("/AdminValidate", middleware.AdminAuth, controllers.AdminValidate)

	r.POST("/Adminlogout", controllers.AdminLogout)

	r.GET("/findall", middleware.AdminAuth, controllers.FindAll)

	r.POST("/finduser", middleware.AdminAuth, controllers.FindUser)

	r.DELETE("/deleteuser", middleware.AdminAuth, controllers.DeleteUser)

	r.Run() // listen and serve on 0.0.0.0:8080 start the server on localhost:8080

}
