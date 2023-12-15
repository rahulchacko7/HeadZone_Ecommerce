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

	TotalRevenue() (models.DashboardRevenue, error)
	DashBoardOrder() (models.DashboardOrder, error)
	AmountDetails() (models.DashboardAmount, error)
	DashBoardUserDetails() (models.DashBoardUser, error)
	DashBoardProductDetails() (models.DashBoardProduct, error)

	SalesByYear(yearInt int) ([]models.OrderDetailsAdmin, error)
	SalesByMonth(monthInt int) ([]models.OrderDetailsAdmin, error)
	SalesByDay(dayInt int) ([]models.OrderDetailsAdmin, error)
}
