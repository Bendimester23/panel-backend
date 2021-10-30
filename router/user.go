package router

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/ultimatepanel2000/backend/db"
)

func InitUser(r *gin.RouterGroup) {
	r.Use(NeedsAuth)
	r.GET("/me", func(c *gin.Context) {
		if !c.MustGet("auth").(bool) {
			return
		}
		user, err := db.DB.User.FindFirst(
			db.User.ID.Equals(c.GetString("user_id")),
		).With(
			db.User.Servers.Fetch(),
		).Exec(ctx)

		if err != nil {
			if errors.Is(err, db.ErrNotFound) {
				c.JSON(404, gin.H{
					"error": "no user found",
				})
				return
			}
			c.JSON(500, gin.H{
				"error": "db error",
			})
			return
		}

		c.JSON(200, user)
	})
}
