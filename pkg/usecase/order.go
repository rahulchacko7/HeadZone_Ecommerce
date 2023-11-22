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
	orderRepository interfaces.OrderRepository
	userUseCase     services.UserUseCase
}

func NewOrderUseCase(repo interfaces.OrderRepository, userUseCase services.UserUseCase) services.OrderUseCase {
	return &orderUseCase{
		orderRepository: repo,
		userUseCase:     userUseCase,
	}
}
func (i *orderUseCase) OrderItemsFromCart(userID, addressID, paymentID int) error {
	cart, err := i.userUseCase.GetCart(userID)
	if err != nil {
		return err
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

	for _, v := range cart.Data {
		if err := i.userUseCase.RemoveFromCart(cart.ID, v.ID); err != nil {
			return err
		}
	}

	return nil
}

func (i *orderUseCase) GetOrders(orderId int) (domain.OrderResponse, error) {

	orders, err := i.orderRepository.GetOrders(orderId)
	if err != nil {
		return domain.OrderResponse{}, err
	}
	return orders, nil
}

func (i *orderUseCase) CancelOrder(orderID int) error {
	orderStatus, err := i.orderRepository.CheckOrderStatusByID(orderID)
	if err != nil {
		return err
	}

	if orderStatus != "PENDING" {
		return errors.New("order cannot be canceled, kindly return the product if accidentally booked")
	}

	err = i.orderRepository.CancelOrder(orderID)
	if err != nil {
		return err
	}

	return nil
}

func (i *orderUseCase) GetAllOrders(userId, page, pageSize int) ([]models.OrderDetails, error) {
	fmt.Println("dddd", userId)
	allorder, err := i.orderRepository.GetAllOrders(userId, page, pageSize)
	if err != nil {
		return []models.OrderDetails{}, err
	}
	return allorder, nil
}
