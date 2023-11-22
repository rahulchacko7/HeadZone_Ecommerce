package interfaces

import (
	"HeadZone/pkg/domain"
	"HeadZone/pkg/utils/models"
)

type OrderRepository interface {
	OrderItems(userid, addressid, paymentid int, total float64) (int, error)
	AddOrderProducts(order_id int, cart []models.GetCart) error
	GetOrders(orderId int) (domain.OrderResponse, error)
}
