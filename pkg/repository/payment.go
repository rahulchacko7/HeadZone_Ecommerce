package repository

import (
	"HeadZone/pkg/repository/interfaces"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type paymentRepositoryImpl struct {
	DB *gorm.DB
}

func NewPaymentRepository(DB *gorm.DB) interfaces.PaymentRepository {
	return &paymentRepositoryImpl{
		DB: DB,
	}
}

// --------------------------------------- add payment details ----------------------------------------- \\

func (repo *paymentRepositoryImpl) AddRazorPayDetails(orderId int, razorPayId string) error {
	query := `
	insert into payments (order_id,razer_id) values($1,$2) 
	`
	if err := repo.DB.Exec(query, orderId, razorPayId).Error; err != nil {
		err = errors.New("error in inserting values to razor pay data table" + err.Error())
		return err
	}
	return nil
}

// ---------------------------------------- update payment details ------------------------------------------- \\

func (repo *paymentRepositoryImpl) UpdatePaymentDetails(orderId string, paymentId string) error {
	fmt.Println("razerId,paymetnId", orderId, paymentId)
	if err := repo.DB.Exec("update payments set payment = $1 where razer_id = $2", paymentId, orderId).Error; err != nil {
		err = errors.New("error in updating the razer pay table " + err.Error())
		return err
	}
	return nil
}

// ------------------------------------------- check payment status ----------------------------------- \\

func (repo *paymentRepositoryImpl) GetPaymentStatus(orderId string) (bool, error) {
	var paymentStatus string
	err := repo.DB.Raw("select payment_status from orders where id = $1", orderId).Scan(&paymentStatus).Error
	if err != nil {
		return false, err
	}

	// Check if payment status is "PAID"
	isPaid := paymentStatus == "PAID"
	fmt.Println("Is payment status PAID?", isPaid)
	return isPaid, nil
}

// -------------------------------------------- update payment status ---------------------------------- \\

func (repo *paymentRepositoryImpl) UpdatePaymentStatus(status bool, orderId string) error {
	var paymentStatus string
	if status {
		paymentStatus = "PAID"
	} else {
		paymentStatus = "NOT PAID"
	}

	query := `
		UPDATE orders SET payment_status = $1, order_status = 'SHIPPED' WHERE id = $2 
	`
	if err := repo.DB.Exec(query, paymentStatus, orderId).Error; err != nil {
		err = errors.New("error in updating orders payment status: " + err.Error())
		return err
	}
	return nil
}
