package interfaces

type PaymentRepository interface {
	AddRazorPayDetails(int, string) error
	UpdatePaymentDetails(orderId string, paymentId string) error
	GetPaymentStatus(orderId string) (bool, error)
	UpdatePaymentStatus(status bool, orderId string) error
}
