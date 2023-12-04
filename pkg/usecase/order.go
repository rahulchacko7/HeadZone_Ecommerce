package usecase

import (
	domain "HeadZone/pkg/domain"
	interfaces "HeadZone/pkg/repository/interfaces"
	services "HeadZone/pkg/usecase/interfaces"
	"HeadZone/pkg/utils/models"
	"errors"
	"fmt"
)

type orderUseCase struct {
	orderRepository  interfaces.OrderRepository
	userUseCase      services.UserUseCase
	walletRepository interfaces.WalletRepository
	cartRepo         interfaces.CartRepository
}

func NewOrderUseCase(repo interfaces.OrderRepository, userUseCase services.UserUseCase, walletRepo interfaces.WalletRepository, cartRepo interfaces.CartRepository) services.OrderUseCase {
	return &orderUseCase{
		orderRepository:  repo,
		userUseCase:      userUseCase,
		walletRepository: walletRepo,
		cartRepo:         cartRepo,
	}
}
func (i *orderUseCase) OrderItemsFromCart(userID, addressID, paymentID int) error {

	if userID <= 0 || addressID <= 0 || paymentID < 0 {
		return errors.New("enter a valid number")
	}

	cart, err := i.userUseCase.GetCart(userID)
	if err != nil {
		return err
	}

	exist, err := i.cartRepo.CheckCart(userID)

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
	orderID, err := i.orderRepository.OrderItems(userID, addressID, paymentID, total)
	if err != nil {
		return err
	}
	fmt.Println("orderid:......", orderID)
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

		// For admin approval, change SHIPPED to DELIVERED
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
