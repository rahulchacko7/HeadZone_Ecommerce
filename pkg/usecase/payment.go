package usecase

import (
	usecase "HeadZone/pkg/repository/interfaces"
	"HeadZone/pkg/usecase/interfaces"
	"HeadZone/pkg/utils/models"
	"errors"
	"fmt"

	"github.com/razorpay/razorpay-go"
)

type paymentUsecaseImpl struct {
	paymentRepo     usecase.PaymentRepository
	orderRepository usecase.OrderRepository
}

func NewPaymentUseCase(repo usecase.OrderRepository, payment usecase.PaymentRepository) interfaces.PaymentUseCase {
	return &paymentUsecaseImpl{
		orderRepository: repo,
		paymentRepo:     payment,
	}
}

// ---------------------------------------- make payment through razor pay --------------------------------------- \\

func (repo *paymentUsecaseImpl) MakePaymentRazorpay(orderId, userId int) (models.CombinedOrderDetails, string, error) {

	fmt.Println("Order iddddddd inside usecase", orderId)

	order, err := repo.orderRepository.GetOrder(orderId)
	if err != nil {
		err = errors.New("error in getting order details through order id" + err.Error())
		return models.CombinedOrderDetails{}, "", err
	}

	client := razorpay.NewClient("rzp_test_dch7hG3p7YJMuI", "S9vb9BEvMMmH94veyS62OZ11")

	fmt.Println("order amount", order.FinalPrice)
	data := map[string]interface{}{
		"amount":   int(order.FinalPrice) * 100,
		"currency": "INR",
		"receipt":  "some_receipt_id",
	}

	body, err := client.Order.Create(data, nil)
	if err != nil {
		return models.CombinedOrderDetails{}, "", nil
	}
	fmt.Println("body usecase", body)
	razorPayOrderId := body["id"].(string)

	err = repo.paymentRepo.AddRazorPayDetails(orderId, razorPayOrderId)
	if err != nil {
		return models.CombinedOrderDetails{}, "", err
	}
	body2, err := repo.orderRepository.GetDetailedOrderThroughId(int(order.ID))
	if err != nil {
		return models.CombinedOrderDetails{}, "", err
	}
	fmt.Println("body 2 usecase", body2.OrderId)

	return body2, razorPayOrderId, nil
}

// ------------------------------------------------- verify payment razor pay ------------------------------------ \\

func (repo *paymentUsecaseImpl) SavePaymentDetails(paymentId, razorId, orderId string) error {

	status, err := repo.paymentRepo.GetPaymentStatus(orderId)
	if err != nil {
		return err
	}
	fmt.Println("status", status)
	if !status {
		err = repo.paymentRepo.UpdatePaymentDetails(razorId, paymentId)
		if err != nil {
			return err
		}

		err = repo.paymentRepo.UpdatePaymentStatus(true, orderId)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("already paid")

}
