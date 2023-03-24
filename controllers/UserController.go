package controllers

import (
	"fmt"
	"net/http"
	"os"
	"tess/initializers"
	"tess/models"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func UserSignup(c *gin.Context) {
	//get the email and the pass off req body
	var body struct {
		Name     string
		Email    string
		Password string
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to read body",
		})
		return
	}

	fmt.Println(body)

	//hash the password

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to hash password",
		})

		return
	}

	user := models.User{Name: body.Name, Email: body.Email, Password: string(hash)}

	result := initializers.DB.Create(&user) // pass pointer of data to Create

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to create user",
		})

		return
	}

	//responds

	c.JSON(http.StatusOK, gin.H{})

}

func UserLogin(c *gin.Context) {
	//get the email and password of req body
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to read body",
		})
		return
	}

	//look up request user
	var user models.User
	initializers.DB.Find(&user, "email= ?", body.Email)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalied user or password",
		})
		return
	}
	//compare sent in pass with hash pass

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "inavied password",
		})
		return
	}

	//generate jwt
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Second * 2 * 30).Unix(),
	})
	//Sign and get the complete encoded token as a string useing secrete

	sc := os.Getenv("SECRET")

	tokenString, err := token.SignedString([]byte(sc))

	if err != nil {
		fmt.Println("test4", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to create token",
		})
		return
	}

	//sent it backBrototype

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 2*30, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"log in": "logged in ",
	})

}

func Validate(c *gin.Context) {
	user, _ := c.Get("user")

	//user.(models.User).

	c.JSON(http.StatusOK, gin.H{
		"massage": user,
	})
}

func UserLogout(c *gin.Context) {

	//c.Header("Cache-Control", "no-cache,no-store,must-revalidate")

	tokenString, err := c.Cookie("Authorization")

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "not logged in",
		})
		return
	}

	// c.SetSameSite(http.SameSiteLaxMode)

	c.SetCookie("Authorization", tokenString, -1, "", "", false, true)

	c.JSON(http.StatusSeeOther, gin.H{
		"message": "logouted successfully ",
	})

	c.Redirect(http.StatusFound, "/")
}

func EditUser(c *gin.Context) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to hash password",
		})

		return
	}

	user := models.User{Password: string(hash)}

	result := initializers.DB.Model(&user).Where("email=?", body.Email).Update("password", string(hash))

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to update password",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "successfully changed password",
	})

}
