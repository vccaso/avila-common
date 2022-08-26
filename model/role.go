package model

import (
	"encoding/json"
	"io"

	"github.com/go-playground/validator/v10"
)

// swagger:model
type Role struct {
	Id        int64  `json:"id"`
	Name      string `json:"name" validate:"required"`
	Isactive  bool   `json:"isactive"`
	Isdeleted bool   `json:"isdeleted"`
}

// swagger:model
type Roles []*Role

// json converters

func (u *Roles) ToJson(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(u)
}

func (u *Role) ToJson(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(u)
}

func (u *Role) FromJson(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(u)
}

// validators

func (u *Role) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}
