package router

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

var expiration_time = (time.Hour * 72).Seconds()

func NeedsAuth(c *gin.Context) {
	tokenString := c.Request.Header.Get("Authorization")

	token, err := jwt.Parse(strings.TrimSpace(tokenString), func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(tokenSecret), nil
	})

	if err != nil {
		c.JSON(401, gin.H{
			"error": "invalid token",
		})
		c.Set("auth", false)
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if int64(claims["date"].(float64))+int64(expiration_time) < time.Now().Unix() {
			c.JSON(401, gin.H{
				"error": "expired token",
			})
			c.Set("auth", false)
			return
		}

		c.Set("auth", true)
		c.Set("user_id", claims["user_id"].(string))
		c.Set("is_admin", claims["is_admin"].(bool))
		return
	} else {
		c.JSON(401, gin.H{
			"error": "invalid token",
		})
		c.Set("auth", false)
		return
	}
}

func NeedsAdmin(c *gin.Context) {
	if c.MustGet("is_admin").(bool) {
		return
	}

	c.JSON(403, gin.H{
		"error": "insufficient role",
	})
	c.Set("auth", false)
}
