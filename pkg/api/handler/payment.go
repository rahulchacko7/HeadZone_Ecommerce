package handler

import (
	"HeadZone/pkg/usecase/interfaces"
	"HeadZone/pkg/utils/response"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	payment interfaces.PaymentUseCase
}

func NewPaymentHandler(useCase interfaces.PaymentUseCase) *PaymentHandler {
	return &PaymentHandler{
		payment: useCase,
	}
}

func (handler *PaymentHandler) MakePaymentRazorpay(c *gin.Context) {
	userId := c.Query("user_id")
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {

		errRes := response.ClientResponse(http.StatusBadRequest, "error", nil, errors.New("error in converting string to int userid"+err.Error()))
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	orderId := c.Query("order_id")
	orderIdInt, err := strconv.Atoi(orderId)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "error", nil, errors.New("error in converting string to int orderid"+err.Error()))
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	fmt.Println("Order iddddddd", orderIdInt)
	body, razorId, err := handler.payment.MakePaymentRazorpay(orderIdInt, userIdInt)
	fmt.Println(" inside Order iddddddd", body.OrderId)

	fmt.Println("body in handler", body)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "error", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	fmt.Println("body next", body.FinalPrice, razorId, userId, body.OrderId, body.Name, body.FinalPrice)

	c.HTML(http.StatusOK, "index.html", gin.H{
		"final_price": body.FinalPrice * 100,
		"razor_id":    razorId,
		"user_id":     userId,
		"order_id":    body.OrderId,
		"user_name":   body.Name,
		"total":       int(body.FinalPrice),
	})
}

func (handler *PaymentHandler) VerifyPayment(c *gin.Context) {
	orderId := c.Query("order_id")
	paymentId := c.Query("payment_id")
	razorId := c.Query("razor_id")

	fmt.Println("details order, payment, razor", orderId)
	fmt.Println("details order, payment, razor", paymentId)
	fmt.Println("details order, payment, razor", razorId)

	if err := handler.payment.SavePaymentDetails(paymentId, razorId, orderId); err != nil {
		errorRes := response.ClientResponse(http.StatusInternalServerError, "could not update payment details", nil, err.Error())
		c.JSON(http.StatusOK, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully updated payment details", nil, nil)
	c.JSON(http.StatusOK, successRes)
}
