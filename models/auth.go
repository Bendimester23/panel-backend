package models

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type Validatable interface {
	Validate() error
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type SetPasswordRequest = LoginRequest

func (m LoginRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Email, validation.Required, is.Email),
		validation.Field(&m.Password, validation.Required, validation.Length(4, 35)),
	)
}

func (m RegisterRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Email, validation.Required, is.Email),
		validation.Field(&m.Password, validation.Required, validation.Length(4, 35)),
		validation.Field(&m.Username, validation.Required, validation.Length(4, 30)),
	)
}
