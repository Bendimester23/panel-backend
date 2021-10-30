package utils

import (
	"github.com/gin-gonic/gin"
)

func Validate(c *gin.Context, target interface{}, t func(t interface{}) error) bool {
	if err := c.BindJSON(target); err != nil {
		c.JSON(400, gin.H{
			"error": "bad request",
		})
		return true
	}

	if err := t(target); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return true
	}

	return false
}
