package interfaces

import (
	"HeadZone/pkg/domain"
	"HeadZone/pkg/utils/models"
)

type AdminRepository interface {
	LoginHandler(adminDetails models.AdminLogin) (domain.Admin, error)
	GetUserByID(id string) (domain.Users, error)
	UpdateBlockUserByID(user domain.Users) error
	GetUsers(page int) ([]models.UserDetailsAtAdmin, error)

	NewPaymentMethod(string) error
	ListPaymentMethods() ([]domain.PaymentMethod, error)
	GetPaymentMethod() ([]models.PaymentMethodResponse, error)
	CheckIfPaymentMethodAlreadyExists(payment string) (bool, error)
	DeletePaymentMethod(id int) error
}
