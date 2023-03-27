package helper

import (
	"net/http"
	"os"
	"tess/initializers"
	"tess/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func AdminCreate(c *gin.Context) {

	//1.get the req from env
	AE := os.Getenv("AdminEmail")
	AP := os.Getenv("AdminPass")

	//2.Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(AP), 10)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to hash password",
		})
		return
	}
	//3.create admin table in db
	admin := models.Admin{Email: AE, Password: string(hash)}

	result := initializers.DB.Create(&admin) // pass pointer of data to Create

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to create ADMIN",
		})
		return
	}

	//4.4esponds

	c.JSON(http.StatusOK, gin.H{
		"Message": "New admin is added",
	})
}
