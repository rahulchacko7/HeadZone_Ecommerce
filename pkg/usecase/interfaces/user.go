package interfaces

import (
	"HeadZone/pkg/domain"
	"HeadZone/pkg/utils/models"
)

type UserUseCase interface {
	UserSignUp(models.UserDetails) (models.TokenUsers, error)
	LoginHandler(models.UserLogin) (models.TokenUsers, error)
	GetUserDetails(id int) (models.UserDetailsResponse, error)
	GetCart(id int) (models.GetCartResponse, error)
	AddAddress(id int, address models.AddAddress) error
	GetAddresses(id int) ([]domain.Address, error)

	EditDetails(int, models.EditDetailsResponse) (models.EditDetailsResponse, error)

	ChangePassword(id int, old string, password string, repassword string) error

	RemoveFromCart(cart, inventory int) error
	UpdateQuantity(id, inv_id, qty int) error
}
