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
	AE := os.Getenv("AdminEmail")
	AP := os.Getenv("AdminPass")

	hash, err := bcrypt.GenerateFromPassword([]byte(AP), 10)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to hash password",
		})
		return
	}
	admin := models.Admin{Email: AE, Password: string(hash)}

	result := initializers.DB.Create(&admin) // pass pointer of data to Create

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to create ADMIN",
		})
		return
	}

	//responds

	c.JSON(http.StatusOK, gin.H{})
}
