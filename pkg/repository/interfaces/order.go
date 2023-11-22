package interfaces

import (
	"HeadZone/pkg/domain"
	"HeadZone/pkg/utils/models"
)

type OrderRepository interface {
	OrderItems(userid, addressid, paymentid int, total float64) (int, error)
	AddOrderProducts(order_id int, cart []models.GetCart) error
	GetOrders(orderId int) (domain.OrderResponse, error)
	// GetOrderDetails(userId int, page int, count int) (models.AllOrderResponse, error)
	CheckOrderStatusByID(id int) (string, error)
	CancelOrder(id int) error
}
