package model

import (
	"database/sql"
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

// function to scan Role object from database single row
func (u *Role) FromRow(row *sql.Rows) error {

	err := row.Scan(&u.Id, &u.Name, &u.Isactive, &u.Isdeleted)
	return err
}

// function to scan Role object from database rows
func (u *Roles) FromRows(rows *sql.Rows) error {

	for rows.Next() {
		r := new(Role)
		err := r.FromRow(rows)
		if err != nil {
			return err
		}
		*u = append(*u, r)
	}
	return nil
}

// validators

func (u *Role) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}
