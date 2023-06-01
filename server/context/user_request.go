package context

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// User Register Request
type UserRequest struct {
	UserName    string `json:"username" validate:"required"`
	UserContact string `json:"usercontact" validate:"required,len=10"`
}

func (u UserRequest) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.UserName,
			validation.Required.Error("Username is required"),
			validation.Match(regexp.MustCompile("^[a-zA-Z]+$")).Error("Username must contain only alphabetic characters"),
		),
	)
}

// User Login Request
type UserLogin struct {
	UserContact string `json:"usercontact" validate:"required,len=10"`
}

// Otp verification Request
type VerifyOtp struct {
	UserContact string `json:"usercontact" validate:"required,len=10"`
	Otp         string `json:"otp" validate:"required,len=6"`
}

// Edit User Details
type EditUser struct {
	Gender   string `json:"gender"`
	UserName string `json:"username"`
}

func (u EditUser) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.UserName,
			validation.Required.Error("Username is required"),
			validation.Match(regexp.MustCompile("^[a-zA-Z]+$")).Error("Username must contain only alphabetic characters"),
		),
		validation.Field(&u.Gender,
			validation.Required.Error("Gender is required"),
			validation.In("male", "female").Error("Gender must be either 'male' or 'female'"),
		),
	)
}

// Logout User
type LogoutUser struct {
	UserId string `json:"userId" validate:"required"`
}
