package controllers

import (
	"fmt"
	"net/http"
	"os"
	"tess/initializers"
	"tess/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

//---------------------------------------------------AdminLogin--------------------------------------------------

func Loginadmin(c *gin.Context) {

	//1.helper.AdminCreate(c)

	var body struct {
		Email    string
		Password string
	}

	//2.bind req body to a struct or map common in json

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to read body",
		})
		return
	}

	//3.check the user is registered or not

	var admin models.Admin
	initializers.DB.Find(&admin, "email= ?", body.Email)

	//4.if no user is available

	if admin.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid Admin ",
		})
		return
	}
	//5.compare sent in pass with hash pass

	err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(body.Password))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid password",
		})
		return
	}

	//6.generate jwt
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": admin.ID,
		"exp": time.Now().Add(time.Second * 2 * 30).Unix(),
	})

	//7.Sign and get the complete encoded token as a string using secrete

	Ab := os.Getenv("SECRET")

	tokenString, err := token.SignedString([]byte(Ab))

	if err != nil {
		fmt.Println("test4", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to create token",
		})
		return
	}
	//8.sendind it back to user cookie

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Auth", tokenString, 2*30, "", "", false, true)

	c.JSON(200, gin.H{"message": "Admin Loged successfully"})
}

func AdminValidate(c *gin.Context) {
	admin, _ := c.Get("admin")

	//admin.(models.admin).

	c.JSON(http.StatusOK, gin.H{
		"massage": admin,
	})
}

//------------------------------------------------AdminLogout---------------------------------------------------

func AdminLogout(c *gin.Context) {

	tokenString, err := c.Cookie("Auth")

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "not logged in",
		})
		return
	}

	// c.SetSameSite(http.SameSiteLaxMode)

	c.SetCookie("Auth", tokenString, -1, "", "", false, true)

	c.JSON(http.StatusSeeOther, gin.H{
		"message": "logouted successfully",
	})

	c.Redirect(http.StatusFound, "/")
}

func FindAll(c *gin.Context) {
	var users []models.User
	result := initializers.DB.Find(&users)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No users found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": users,
	})

}

//-------------------------------------------------Find_Users-------------------------------------------------

func FindUser(c *gin.Context) {

	//geting username and email

	var body struct {
		Email string `json:"email"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to read body",
		})
		return
	}

	fmt.Println("bind", body.Email)

	var user models.User
	initializers.DB.Where(" email = ?", body.Email).Find(&user)
	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to find the user",
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": user,
		})
	}

}

//---------------------------------------------------DELETE_USER---------------------------------------------

func DeleteUser(c *gin.Context) {
	var body struct {
		Name  string
		Email string
	}

	err := c.Bind(&body)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "failed to read body",
		})
		return
	}

	// var user models.User
	// initializers.DB.Model(&user).Where("name = ? AND email= ?", body.Name, body.Email).First(&user, 1).Delete(&user)
	initializers.DB.Exec("DELETE FROM users WHERE name = $1 AND email = $2", body.Name, body.Email)
	c.JSON(http.StatusOK, gin.H{
		"message": "user deleted",
	})

}

// --------------------------------------------------Edit_User-------------------------------------------------
func EditUser(c *gin.Context) {
	var body struct {
		Email    string
		Password string
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
