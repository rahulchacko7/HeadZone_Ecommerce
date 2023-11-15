package interfaces

import (
	"HeadZone/pkg/domain"
	"HeadZone/pkg/utils/models"
)

type UserUseCase interface {
	UserSignUp(models.UserDetails) (models.TokenUsers, error)
	LoginHandler(models.UserLogin) (models.TokenUsers, error)
	GetUserDetails(id int) (models.UserDetailsResponse, error)
	AddAddress(id int, address models.AddAddress) error
	GetAddresses(id int) ([]domain.Address, error)

	EditName(id int, name string) error
	EditEmail(id int, email string) error
	EditPhone(id int, phone string) error

	ChangePassword(id int, old string, password string, repassword string) error
}
