package interfaces

import "HeadZone/pkg/domain"

type OrderUseCase interface {
	OrderItemsFromCart(userid int, addressid int, paymentid int) error
	GetOrders(orderId int) (domain.OrderResponse, error)
	GetAllOrders(id int)
}
