package usermodel

type UserModel struct {
	Id          string `json:"id"`
	PhoneNumber string `json:"phone_number" validate:"required,len=10,numeric"`
	OTP         string `json:"otp" validate:"required,len=6,numeric"`
	Token       string `json:"token"`
}
