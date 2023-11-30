package interfaces

import "HeadZone/pkg/utils/models"

type PaymentUseCase interface {
	MakePaymentRazorpay(orderId, userId int) (models.CombinedOrderDetails, string, error)
	SavePaymentDetails(paymentId, razorId, orderId string) error
}
