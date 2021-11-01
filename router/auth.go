package router

import (
	"context"
	"errors"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/golang-jwt/jwt"
	"github.com/ultimatepanel2000/backend/db"
	"github.com/ultimatepanel2000/backend/models"
	"github.com/ultimatepanel2000/backend/utils"
	"golang.org/x/crypto/bcrypt"
)

var ctx = context.Background()

var tokenSecret = []byte("dhfhfgh")

func InitAuth(r *gin.RouterGroup) {
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "running",
		})
	})

	r.PUT("/register", func(c *gin.Context) {
		var data = models.RegisterRequest{}

		if utils.Validate(c, &data, func(t interface{}) error {
			return t.(*models.RegisterRequest).Validate()
		}) {
			return
		}

		_, err := db.DB.User.FindFirst(
			db.User.Email.Equals(data.Email),
		).Exec(ctx)

		if !errors.Is(err, db.ErrNotFound) {
			c.JSON(409, gin.H{
				"error": "email already in use",
			})
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.MinCost)

		if handleError(err, c, "error creating hash") {
			return
		}

		user, e := db.DB.User.CreateOne(
			db.User.Name.Set(data.Username),
			db.User.Password.Set(string(hash)),
			db.User.Email.Set(data.Email),
		).Exec(ctx)

		if handleError(e, c, "db error") {
			return
		}

		token, err := newToken(user.ID, user.IsAdmin)

		if handleError(err, c, "error creating token") {
			return
		}

		c.JSON(200, gin.H{
			"status": "success",
			"token":  token,
		})
	})

	r.POST("/login", func(c *gin.Context) {
		var data = &models.LoginRequest{}

		if utils.Validate(c, data, func(t interface{}) error {
			return t.(*models.LoginRequest).Validate()
		}) {
			return
		}

		user, err := db.DB.User.FindFirst(
			db.User.Email.Equals(data.Email),
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

		if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password)) != nil {
			c.JSON(403, gin.H{
				"error": "invalid credentials",
			})
			return
		}

		token, err := newToken(user.ID, user.IsAdmin)

		if handleError(err, c, "error creating token") {
			return
		}

		c.JSON(200, gin.H{
			"token": token,
		})
	})

	r.Use(NeedsAuth).GET("/refresh", func(c *gin.Context) {
		if !c.MustGet("auth").(bool) {
			return
		}

		user, err := db.DB.User.FindFirst(
			db.User.ID.Equals(c.GetString("user_id")),
		).Exec(ctx)

		if err != nil {
			if errors.Is(err, db.ErrNotFound) {
				c.JSON(404, gin.H{
					"error": "user not found",
				})
				return
			}
			c.JSON(500, gin.H{
				"error": "db error",
			})
		}

		token, err := newToken(c.GetString("user_id"), user.IsAdmin)

		if handleError(err, c, "token creation error") {
			return
		}

		c.JSON(200, gin.H{
			"token": token,
		})
	})

	r.PATCH("/password", func(c *gin.Context) {
		p, _ := ioutil.ReadAll(c.Request.Body)
		password := string(p)

		if validation.Validate(&password, validation.Required, validation.Length(4, 35)) != nil {
			c.JSON(400, gin.H{
				"error": "too short or too long password",
			})
			return
		}

		c.JSON(200, gin.H{
			"status": "success",
		})
	})

	r.Use(NeedsAdmin).PATCH("/setpassword", func(c *gin.Context) {
		var data = &models.SetPasswordRequest{}

		if utils.Validate(c, &data, func(t interface{}) error {
			return t.(*models.SetPasswordRequest).Validate()
		}) {
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)

		if handleError(err, c, "error generating hash") {
			return
		}

		_, err = db.DB.User.FindMany(
			db.User.Email.Equals(data.Email),
		).Update(
			db.User.Password.Set(string(hash)),
		).Exec(ctx)

		if handleError(err, c, "db error") {
			return
		}
	})

	r.DELETE("/:email", func(c *gin.Context) {
		email := c.Param("email")
		if err := validation.Validate(&email, validation.Required, is.Email); err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}

		_, err := db.DB.User.FindMany(
			db.User.Email.Equals(email),
		).Delete().Exec(ctx)

		if handleError(err, c, "db error") {
			return
		}

		c.JSON(200, gin.H{
			"status": "success",
		})
	})
}

func newToken(id string, isAdmin bool) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"date":     time.Now().Unix(),
		"user_id":  id,
		"is_admin": isAdmin,
	})

	return t.SignedString(tokenSecret)
}

func handleError(err error, c *gin.Context, erro string) bool {
	if err != nil {
		c.JSON(500, gin.H{
			"error": erro,
		})
		return true
	}

	return false
}
