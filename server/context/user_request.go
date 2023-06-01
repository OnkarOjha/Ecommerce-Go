package context

// User Register Request
type UserRequest struct {
	UserName    string `json:"username" validate:"required"`
	UserContact string `json:"usercontact" validate:"required,len=10"`
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
	Gender   string `json:"gender"  validate:"oneof=male female"`
	UserName string `json:"username"`
}

// Logout User
type LogoutUser struct {
	UserId string `json:"userId" validate:"required"`
}
