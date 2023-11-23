package interfaces

import (
	"HeadZone/pkg/domain"
	"HeadZone/pkg/utils/models"
)

type OrderRepository interface {
	OrderItems(userid, addressid, paymentid int, total float64) (int, error)
	AddOrderProducts(order_id int, cart []models.GetCart) error
	GetOrders(orderId int) (domain.OrderResponse, error)
	CheckOrderStatusByID(id int) (string, error)
	CancelOrder(id int) error
	GetAllOrders(userId, page, pageSize int) ([]models.OrderDetails, error)
	GetOrderDetailsBrief(page int) ([]models.CombinedOrderDetails, error)
	CheckOrderStatus(orderId string) (bool, error)
	GetShipmentStatus(orderId string) (string, error)
	ApproveOrder(orderId string) error
}
