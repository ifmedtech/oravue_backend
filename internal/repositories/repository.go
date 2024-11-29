package repositories

type Repository interface {
	CreateUser(phoneNumber string, otp string) (string, error)
}
