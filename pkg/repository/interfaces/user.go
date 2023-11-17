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

	GetCartID(id int) (int, error)
	GetProductsInCart(cart_id int) ([]int, error)
	FindProductNames(inventory_id int) (string, error)
	FindCartQuantity(cart_id, inventory_id int) (int, error)
	FindCategory(inventory_id int) (int, error)

	RemoveFromCart(cart, inventory int) error
	UpdateQuantityAdd(id, inv_id int) error
	UpdateQuantityLess(id, inv_id int) error
}
