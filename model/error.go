package model

import (
	"encoding/json"
	"io"
	"time"

	"github.com/go-playground/validator/v10"
)

// swagger:model
type Error struct {
	// the id for the product
	//
	// required: true
	// min: 1
	Id int64 `json:"id"`
	// Time for the error
	//
	// required: true
	Error_time time.Time `json:"error_time" validate:"required"`
	// application where the error was thrown
	//
	// required: true
	App string `json:"app" validate:"required"`
	// level
	//
	// required: true
	Level string `json:"level" validate:"required"`
	// Error message
	//
	// required: true
	Message string `json:"message" validate:"required"`
	// UUID with the gateway session
	//
	// required: true
	Gateway_session string `json:"gateway_session" validate:"required"`
}

// swagger:model
type Errors []*Error

// json converters

func (u *Errors) ToJson(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(u)
}

func (u *Error) ToJson(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(u)
}

func (u *Error) FromJson(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(u)
}

// validators

func (u *Error) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}
