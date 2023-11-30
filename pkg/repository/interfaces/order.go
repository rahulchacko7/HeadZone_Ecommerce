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
	CheckOrdersStatusByID(id string) (string, error)
	GetShipmentStatus(orderId string) (string, error)
	ApproveOrder(orderId string) error
	ChangeOrderStatus(orderID, status string) error
	GetShipmentsStatus(orderID string) (string, error)
	ReturnOrder(shipmentStatus string, orderID string) error
	ReduceInventoryQuantity(productName string, quantity int) error
	GetOrderDetailsByOrderId(orderID string) (models.CombinedOrderDetails, error)
	AddRazorPayDetails(orderID string, razorPayOrderID string) error
	GetOrder(int) (domain.Order, error)
	GetOrdersDetailsByOrderId(orderID int) (models.CombinedOrderDetails, error)
	PaymentMethodID(orderID int) (int, error)
	PaymentAlreadyPaid(orderID int) (bool, error)
	GetDetailedOrderThroughId(orderId int) (models.CombinedOrderDetails, error)
}
