package model

import (
	"encoding/json"
	"io"
	"time"

	"github.com/google/uuid"

	"github.com/go-playground/validator/v10"
)

// swagger:model
type User struct {
	Session uuid.UUID `json:"session"`
	// the id for the product
	//
	// required: false
	// min: 1
	Id            int64     `json:"id"`
	Email         string    `json:"email" validate:"required,email"`
	Firstname     string    `json:"firstName" validate:"required"`
	Lastname      string    `json:"lastName" validate:"required"`
	Alias         string    `json:"alias" validate:"required"`
	Avatar        string    `json:"avatar" validate:"required"`
	Status        string    `json:"status" validate:"required"`
	ResetPassword bool      `json:"resetPassword"`
	Password      string    `json:"password" validate:"required"`
	CreationDate  time.Time `json:"creationDate"`
	Isactive      bool      `json:"isactive"`
	Isdeleted     bool      `json:"isdeleted"`
}

// swagger:model
type Users []*User

// json converters

func (u *Users) ToJson(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(u)
}

func (u *User) ToJson(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(u)
}

func (u *User) FromJson(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(u)
}

// validators

func (u *User) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}
