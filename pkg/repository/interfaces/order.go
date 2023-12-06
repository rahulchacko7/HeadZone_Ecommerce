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
	CheckPaymentStatus(orderID int) (string, error)
	FindFinalPrice(orderID int) (int, error)
	FindUserID(orderID int) (int, error)
	UpdateOrder(orderID int) ([]models.CombinedOrderDetails, error)
	UpdateReturnedOrder(orderID int) ([]models.CombinedOrderDetails, error)
	CancelOrder(id int) error
	GetAllOrders(userId, page, pageSize int) ([]models.OrderDetails, error)
	GetOrderDetailsBrief(page int) ([]models.CombinedOrderDetails, error)
	CheckOrdersStatusByID(id int) (string, error)
	GetShipmentStatus(orderId int) (string, error)
	ApproveOrder(orderId string) error
	ChangeOrderStatus(orderID int, status string) error
	GetShipmentsStatus(orderID int) (string, error)
	ReturnOrder(shipmentStatus string, orderID int) error
	ReduceInventoryQuantity(productName string, quantity int) error
	GetOrderDetailsByOrderId(orderID string) (models.CombinedOrderDetails, error)
	AddRazorPayDetails(orderID string, razorPayOrderID string) error
	GetOrder(int) (domain.Order, error)
	GetOrdersDetailsByOrderId(orderID int) (models.CombinedOrderDetails, error)
	PaymentMethodID(orderID int) (int, error)
	PaymentAlreadyPaid(orderID int) (bool, error)
	GetDetailedOrderThroughId(orderId int) (models.CombinedOrderDetails, error)
	GetOrderStatus(orderID int) (string, error)
	CheckOrderStatusByOrderId(orderID int) (string, error)
	OrderIdStatus(orderID int) (bool, error)
}
