package interfaces

import (
	"HeadZone/pkg/domain"
	"HeadZone/pkg/utils/models"
)

type UserRepository interface {
	UserSignUp(user models.UserDetails) (models.UserDetailsResponse, error)
	CheckUserAvailability(email string) bool
	FindUserByEmail(user models.UserLogin) (models.UserSignInResponse, error)
	UserBlockStatus(email string) (bool, error)
	GetUserDetails(id int) (models.UserDetailsResponse, error)

	CheckIfFirstAddress(id int) bool
	AddAddress(id int, address models.AddAddress, result bool) error
	GetAddresses(id int) ([]domain.Address, error)

	EditName(id int, name string) error
	EditEmail(id int, email string) error
	EditPhone(id int, phone string) error

	ChangePassword(id int, password string) error
	GetPassword(id int) (string, error)
}
