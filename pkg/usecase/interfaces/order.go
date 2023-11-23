package interfaces

import (
	"HeadZone/pkg/domain"
	"HeadZone/pkg/utils/models"
)

type OrderUseCase interface {
	OrderItemsFromCart(userid int, addressid int, paymentid int) error
	GetOrders(orderId int) (domain.OrderResponse, error)
	GetAllOrders(userId, page, pageSize int) ([]models.OrderDetails, error)
	CancelOrder(orderId int) error
	GetAdminOrders(page int) ([]models.CombinedOrderDetails, error)
}
