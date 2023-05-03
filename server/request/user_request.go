package request

type UserRequest struct {
	UserName    string `json:"username" validate:"required"`
	UserContact string `json:"usercontact" validate:"required,len=10"`
}

type UserLogin struct{
	UserContact string `json:"usercontact" validate:"required,len=10"`
}

type VerifyOtp struct{
	UserContact string `json:"usercontact" validate:"required,len=10"`
	Otp string `json:"otp" validate:"required,len=6"`
}

type GetUser struct{
	UserId string `json:"userId" validate:"required"`
}

type EditUser struct{
	UserId string `json:"userId" validate:"required"`
	Gender string `json:"gender"  validate:"oneof=male female"`
	UserName string `json:"username"`
}


type LogoutUser struct{
	UserId string `json:"userId" validate:"required"`
}