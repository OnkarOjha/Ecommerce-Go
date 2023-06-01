package context

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Vendor Register Request
type VendorRegisterRequest struct {
	GstNumber      string `json:"gstNumber"`
	CompanyName    string `json:"companyName"`
	CompanyContact string `json:"companyContact"`
	Street         string `json:"street"`
	City           string `json:"city"`
	State          string `json:"state"`
	PostalCode     string `json:"postalCode"`
	Country        string `json:"country"`
}

func (a VendorRegisterRequest) ValidateRegister() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.GstNumber, validation.Required),

		validation.Field(&a.CompanyName, validation.Required),

		validation.Field(&a.Street, validation.Required, validation.Length(5, 50)),

		validation.Field(&a.City, validation.Required, validation.Length(5, 15)),

		validation.Field(&a.State, validation.Required, validation.Length(5, 20)),

		validation.Field(&a.PostalCode, validation.Required, validation.Length(1, 6)),

		validation.Field(&a.Country, validation.Required, validation.Length(5, 20)),

		validation.Field(&a.CompanyContact, validation.Required, validation.Match(regexp.MustCompile(`^(?:\+91|0)?[6789]\d{9}$`))),
	)
}

// vendor login request
type VendorLoginRequest struct {
	GstNumber      string `json:"gstNumber" validate:"required"`
	CompanyContact string `json:"companyContact" validate:"required"`
}

func (a VendorLoginRequest) ValidateLogin() error {
	return validation.ValidateStruct(&a,

		validation.Field(&a.CompanyContact, validation.Required, validation.Match(regexp.MustCompile(`^(?:\+91|0)?[6789]\d{9}$`))),
	)
}

// vendor verify otp
type VendorVerifyOtpRequest struct {
	ContactNumber string `json:"contactNumber"`
	Otp           string `json:"otp"`
}

func (a VendorVerifyOtpRequest) ValidateOtp() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.ContactNumber, validation.Required, validation.Match(regexp.MustCompile(`^(?:\+91|0)?[6789]\d{9}$`))),
		validation.Field(&a.Otp, validation.Required, validation.Length(1, 6)),
	)
}

// vendor edit details
type VendorEditDetailsRequest struct {
	GstNumber      string `json:"gstNumber"`
	CompanyName    string `json:"companyName"`
	CompanyContact string `json:"companyContact"`
	Street         string `json:"street"`
	City           string `json:"city"`
	State          string `json:"state"`
	PostalCode     string `json:"postalCode"`
	Country        string `json:"country"`
	Description    string `json:"description"`
	Logo           string `json:"logo"`
	BannerImage    string `json:"bannerImage"`
}

func (a VendorEditDetailsRequest) ValidateEditDetails() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.GstNumber),

		validation.Field(&a.CompanyName),

		validation.Field(&a.Street, validation.Length(5, 50)),

		validation.Field(&a.City, validation.Length(5, 15)),

		validation.Field(&a.State, validation.Length(5, 20)),

		validation.Field(&a.PostalCode, validation.Length(1, 6)),

		validation.Field(&a.Country, validation.Length(5, 20)),

		validation.Field(&a.CompanyContact, validation.Match(regexp.MustCompile(`^(?:\+91|0)?[6789]\d{9}$`))),
	)
}
