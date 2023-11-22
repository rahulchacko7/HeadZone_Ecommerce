package interfaces

import (
	"HeadZone/pkg/domain"
)

type OrderUseCase interface {
	OrderItemsFromCart(userid int, addressid int, paymentid int) error
	GetOrders(orderId int) (domain.OrderResponse, error)
	// GetOrderDetails(userId int, page int, count int) (models.AllOrderResponse, error)
	CancelOrder(orderId int) error
}
