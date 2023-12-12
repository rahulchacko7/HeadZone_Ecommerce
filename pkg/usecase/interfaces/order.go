package interfaces

import (
	"HeadZone/pkg/domain"
	"HeadZone/pkg/utils/models"

	"github.com/jung-kurt/gofpdf"
)

type OrderUseCase interface {
	OrderItemsFromCart(userid int, addressid int, paymentid int, couponid int) error
	GetOrders(orderId int) (domain.OrderResponse, error)
	GetAllOrders(userId, page, pageSize int) ([]models.OrderDetails, error)
	CancelOrder(orderId int) error
	GetAdminOrders(page int) ([]models.CombinedOrderDetails, error)
	OrdersStatus(orderId int) error
	ReturnOrder(orderID int) error
	PaymentMethodID(order_id int) (int, error)
	PrintInvoice(orderIdInt int) (*gofpdf.Fpdf, error)
}
