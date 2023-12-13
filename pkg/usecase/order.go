package usecase

import (
	domain "HeadZone/pkg/domain"
	interfaces "HeadZone/pkg/repository/interfaces"
	services "HeadZone/pkg/usecase/interfaces"
	"HeadZone/pkg/utils/models"
	"errors"
	"fmt"
	"strconv"

	"github.com/jung-kurt/gofpdf"
)

type orderUseCase struct {
	orderRepository  interfaces.OrderRepository
	userUseCase      services.UserUseCase
	walletRepository interfaces.WalletRepository
	cartRepo         interfaces.CartRepository
	couponRepository interfaces.CouponRepository
}

func NewOrderUseCase(repo interfaces.OrderRepository, userUseCase services.UserUseCase, walletRepo interfaces.WalletRepository, cartRepo interfaces.CartRepository, couponRepository interfaces.CouponRepository) services.OrderUseCase {
	return &orderUseCase{
		orderRepository:  repo,
		userUseCase:      userUseCase,
		walletRepository: walletRepo,
		cartRepo:         cartRepo,
		couponRepository: couponRepository,
	}
}
func (i *orderUseCase) OrderItemsFromCart(userID, addressID, paymentID, couponID int) error {

	if userID <= 0 || addressID <= 0 || paymentID < 0 || couponID < 0 {
		return errors.New("enter a valid number")
	}

	cart, err := i.userUseCase.GetCart(userID)
	if err != nil {
		return err
	}

	fmt.Println("cart details at usecase", cart)

	exist, err := i.cartRepo.CheckCart(userID)

	fmt.Println("exist", exist)

	if err != nil {
		return err
	}

	if !exist {
		return errors.New("cart is empty")
	}

	var total float64
	for _, item := range cart.Data {
		if item.Quantity > 0 && item.Price > 0 {
			total += float64(item.Quantity) * float64(item.Price)
		}
	}

	couponvalid, err := i.couponRepository.CheckCouponValid(couponID)
	if err != nil {
		return err
	}
	if !couponvalid {
		return errors.New("this coupon is invalid")
	}

	coupon, err := i.couponRepository.FindCouponPrice(couponID)
	if err != nil {
		return err
	}

	totaldiscount := float64(coupon)

	total = total - totaldiscount

	orderID, err := i.orderRepository.OrderItems(userID, addressID, paymentID, total)
	if err != nil {
		return err
	}
	fmt.Println("orderid use1:......", orderID)
	if err := i.orderRepository.AddOrderProducts(orderID, cart.Data); err != nil {
		return err
	}

	// Update inventory for each product in the cart after a successful order
	for _, v := range cart.Data {
		if err := i.orderRepository.ReduceInventoryQuantity(v.ProductName, v.Quantity); err != nil {
			// Handle error if reducing inventory fails
			return err
		}
	}

	fmt.Println("orderid use2:......", orderID)

	var (
		categoryIds  []int
		productNames []string
		prices       []int
		quantities   []int
		totalPrices  []float64
	)

	fmt.Println("orderid use3:......", orderID)

	for _, item := range cart.Data {
		categoryIds = append(categoryIds, item.Category_id)
		productNames = append(productNames, item.ProductName)
		prices = append(prices, item.Price)
		quantities = append(quantities, item.Quantity)
		totalPrices = append(totalPrices, item.Total)
	}

	fmt.Println("order id at use4 ", orderID)

	err = i.orderRepository.OrderItemsInv(productNames, categoryIds, prices, quantities, totalPrices, userID, orderID)
	if err != nil {
		return errors.New("failed to order items")
	}

	// Remove purchased items from the user's cart
	for _, v := range cart.Data {
		if err := i.userUseCase.RemoveFromCart(cart.ID, v.ID); err != nil {
			return err
		}
	}

	return nil
}

func (i *orderUseCase) GetOrders(orderId int) (domain.OrderResponse, error) {

	if orderId <= 0 {
		return domain.OrderResponse{}, errors.New("enter a valid number")
	}

	orders, err := i.orderRepository.GetOrders(orderId)
	if err != nil {
		return domain.OrderResponse{}, err
	}
	return orders, nil
}

func (i *orderUseCase) CancelOrder(orderID int) error {
	if orderID <= 0 {
		return errors.New("enter a valid number")
	}
	paymentStatus, err := i.orderRepository.CheckPaymentStatus(orderID)
	if err != nil {
		return err
	}

	orderStatus, err := i.orderRepository.CheckOrderStatusByOrderId(orderID)
	if err != nil {
		return err
	}

	price, err := i.orderRepository.FindFinalPrice(orderID)
	if err != nil {
		return err
	}

	userID, err := i.orderRepository.FindUserID(orderID)
	if err != nil {
		return err
	}

	if paymentStatus == "PAID" && orderStatus == "DELIVERED" {
		return errors.New("cannot cancel the item, kindly return it")
	} else if paymentStatus == "PAID" && (orderStatus == "PENDING" || orderStatus == "SHIPPED") {
		// Adding amount back to the user's wallet
		_, errWallet := i.walletRepository.AddToWallet(price, userID)
		if errWallet != nil {
			return errWallet
		}

		// Update order status to CANCELLED if payment is PAID and status is PENDING or SHIPPED
		_, err := i.orderRepository.UpdateOrder(orderID)
		if err != nil {
			return err
		}
	} else {
		// If payment is not PAID or order status is not PENDING or SHIPPED, cancel the order directly
		err = i.orderRepository.CancelOrder(orderID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (i *orderUseCase) GetAllOrders(userId, page, pageSize int) ([]models.OrderDetails, error) {

	if userId <= 0 || page <= 0 || pageSize <= 0 {
		return nil, errors.New("please provide valid input values")
	}

	allorder, err := i.orderRepository.GetAllOrders(userId, page, pageSize)
	if err != nil {
		return []models.OrderDetails{}, err
	}
	return allorder, nil
}

func (i *orderUseCase) GetAdminOrders(page int) ([]models.CombinedOrderDetails, error) {

	if page <= 0 {
		return nil, errors.New("enter a valid number")
	}

	orderDetails, err := i.orderRepository.GetOrderDetailsBrief(page)
	if err != nil {
		return []models.CombinedOrderDetails{}, err
	}
	return orderDetails, nil
}

func (i *orderUseCase) OrdersStatus(orderID int) error {
	if orderID <= 0 {
		return errors.New("enter a valid number")
	}

	// Check if order exists
	orderExists, err := i.orderRepository.OrderIdStatus(orderID)
	if err != nil {
		return err
	}
	if !orderExists {
		return errors.New("no order exists")
	}

	status, err := i.orderRepository.CheckOrdersStatusByID(orderID)
	if err != nil {
		return err
	}

	switch status {
	case "CANCELED", "RETURNED", "DELIVERED":
		return errors.New("cannot approve this order because it's in a processed or canceled state")
	case "PENDING":
		// For admin approval, change PENDING to SHIPPED
		err := i.orderRepository.ChangeOrderStatus(orderID, "SHIPPED")
		if err != nil {
			return err
		}
	case "SHIPPED":
		shipmentStatus, err := i.orderRepository.GetShipmentStatus(orderID)
		if err != nil {
			return err
		}
		if shipmentStatus == "CANCELLED" {
			return errors.New("cannot approve this order because it's cancelled")
		}
		err = i.orderRepository.ChangeOrderStatus(orderID, "DELIVERED")
		if err != nil {
			return err
		}
	}
	return nil
}

func (o *orderUseCase) ReturnOrder(orderID int) error {

	if orderID <= 0 {
		return errors.New("enter a valid number")
	}

	shipmentStatus, err := o.orderRepository.GetOrderStatus(orderID)
	if err != nil {
		return err
	}

	userID, err := o.orderRepository.FindUserID(orderID)
	if err != nil {
		return err
	}

	price, err := o.orderRepository.FindFinalPrice(orderID)
	if err != nil {
		return err
	}

	// Adding amount back to the user's wallet
	_, errWallet := o.walletRepository.AddToWallet(price, userID)
	if errWallet != nil {
		return errWallet
	}

	if shipmentStatus == "DELIVERED" {
		if err := o.orderRepository.ReturnOrder("RETURNED", orderID); err != nil {
			return err
		}
		return nil
	}

	return errors.New("cannot return order")
}

func (or *orderUseCase) PaymentMethodID(order_id int) (int, error) {
	if order_id <= 0 {
		return 0, errors.New("enter a valid number")
	}

	fmt.Println("mmmmmmmmmmmmmmmmm", order_id)
	id, err := or.orderRepository.PaymentMethodID(order_id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (or *orderUseCase) PrintInvoice(orderId int) (*gofpdf.Fpdf, error) {
	order, err := or.orderRepository.GetDetailedOrderThroughId(orderId)
	if err != nil {
		return nil, err
	}

	items, err := or.orderRepository.GetItemsByOrderId(orderId)
	if err != nil {
		return nil, err
	}

	fmt.Println("itemssss", items)

	// Create a new PDF document
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Set font and title
	pdf.SetFont("Arial", "B", 24)
	pdf.SetTextColor(31, 73, 125) // Set text color to a blue shade
	pdf.Cell(0, 20, "Invoice")
	pdf.Ln(20)

	// Customer details section
	pdf.SetFont("Arial", "I", 14)
	pdf.SetTextColor(51, 51, 51) // Set text color to dark gray
	pdf.Cell(0, 10, "Customer Details")
	pdf.Ln(10)
	customerDetails := []string{
		"Name: " + order.Name,
		"House Name: " + order.HouseName,
		"Street: " + order.Street,
		"State: " + order.State,
		"City: " + order.City,
	}
	for _, detail := range customerDetails {
		pdf.Cell(0, 10, detail)
		pdf.Ln(10)
	}
	pdf.Ln(10)

	// Items section headers
	pdf.SetFont("Arial", "B", 16)
	pdf.SetFillColor(217, 217, 217)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(40, 10, "Item", "1", 0, "C", true, 0, "")
	pdf.CellFormat(40, 10, "Price", "1", 0, "C", true, 0, "")
	pdf.CellFormat(40, 10, "Quantity", "1", 0, "C", true, 0, "")
	pdf.CellFormat(40, 10, "Total Price", "1", 0, "C", true, 0, "")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 12)
	pdf.SetFillColor(255, 255, 255) // Set white background for items
	for _, item := range items {
		pdf.CellFormat(40, 10, item.ProductName, "1", 0, "L", true, 0, "")
		pdf.CellFormat(40, 10, "$"+strconv.FormatFloat(item.Price, 'f', 2, 64), "1", 0, "C", true, 0, "")
		pdf.CellFormat(40, 10, strconv.Itoa(item.Quantity), "1", 0, "C", true, 0, "")
		pdf.CellFormat(40, 10, "$"+strconv.FormatFloat(item.Total, 'f', 2, 64), "1", 0, "C", true, 0, "")
		pdf.Ln(10)
	}
	pdf.Ln(10)

	// Total amount section

	var totalPrice float64
	for _, item := range items {
		totalPrice += item.Total
	}

	pdf.SetFont("Arial", "B", 16)
	pdf.SetFillColor(217, 217, 217)
	pdf.CellFormat(120, 10, "Total Price:", "1", 0, "R", true, 0, "")
	pdf.CellFormat(40, 10, "$"+strconv.FormatFloat(totalPrice, 'f', 2, 64), "1", 0, "C", true, 0, "")
	pdf.Ln(10)

	OfferApplied := totalPrice - order.FinalPrice

	fmt.Println("offer Applied", OfferApplied)

	pdf.SetFont("Arial", "B", 16)
	pdf.SetFillColor(217, 217, 217)
	pdf.CellFormat(120, 10, "Offer Applied:", "1", 0, "R", true, 0, "")
	pdf.CellFormat(40, 10, "$"+strconv.FormatFloat(OfferApplied, 'f', 2, 64), "1", 0, "C", true, 0, "")
	pdf.Ln(10)

	pdf.SetFont("Arial", "B", 16)
	pdf.SetFillColor(217, 217, 217)
	pdf.CellFormat(120, 10, "Final Amount:", "1", 0, "R", true, 0, "")
	pdf.CellFormat(40, 10, "$"+strconv.FormatFloat(order.FinalPrice, 'f', 2, 64), "1", 0, "C", true, 0, "")
	pdf.Ln(10)

	return pdf, nil
}
