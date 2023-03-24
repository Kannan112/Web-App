package middleware

import (
	"fmt"
	"net/http"
	"os"
	"tess/initializers"
	"tess/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func AdminAuth(c *gin.Context) {
	fmt.Println("test1")
	//get the cookie off req
	tokenString, err := c.Cookie("Auth")

	
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "you are not logged in",
		})
		return
	}
	//decode/Validate it
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("SECRET")), nil
	})

	//if error to parse the token
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	fmt.Println("test3", token)
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check the exp
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		fmt.Println("test4")
		//find the user with tocken sub
		var admin models.Admin
		initializers.DB.First(&admin, claims["sub"])
		fmt.Println("test5")
		if admin.ID == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		//attach to req
		c.Set("admin", admin)

		//continue

	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	fmt.Println("test6")
}
