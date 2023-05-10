package model

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

// DB model to handle User Addresses
type UserAddresses struct {
	UserId      string `json:"userId"`
	Name        string `json:"name" validate:"required"`
	Street      string `json:"street" validate:"required"`
	City        string `json:"city" validate:"required"`
	State       string `json:"state" validate:"required"`
	PostalCode  string `json:"postalCode" validate:"required"`
	Country     string `json:"country" validate:"required"`
	Phone       string `json:"phone" validate:"required"`
	Email       string `json:"email" validate:"required"`
	AddressType string `json:"addressType"`
}

func (a UserAddresses) Validate() error {
	return validation.ValidateStruct(&a,

		validation.Field(&a.Name, validation.Required),

		validation.Field(&a.Street, validation.Required, validation.Length(5, 50)),

		validation.Field(&a.City, validation.Required, validation.Length(5, 15)),

		validation.Field(&a.State, validation.Required, validation.Length(5, 20)),

		validation.Field(&a.PostalCode, validation.Required, validation.Length(1, 6)),

		validation.Field(&a.Country, validation.Required, validation.Length(5, 20)),

		validation.Field(&a.Email, validation.Required, is.Email),

		validation.Field(&a.Phone, validation.Required, validation.Match(regexp.MustCompile(`^(?:\+91|0)?[6789]\d{9}$`))),

		validation.Field(&a.AddressType, validation.In("default", "work", "home")),
	)
}
