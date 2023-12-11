package handler

import (
	"HeadZone/pkg/usecase/interfaces"
	models "HeadZone/pkg/utils/models"
	response "HeadZone/pkg/utils/response"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderUseCase interfaces.OrderUseCase
}

func NewOrderHandler(useCase interfaces.OrderUseCase) *OrderHandler {
	return &OrderHandler{
		orderUseCase: useCase,
	}
}

func (i *OrderHandler) OrderItemsFromCart(c *gin.Context) {

	userId, _ := c.Get("id")
	UserID, _ := userId.(int)

	var order models.Order
	if err := c.BindJSON(&order); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	if err := i.orderUseCase.OrderItemsFromCart(UserID, order.AddressID, order.PaymentMethodID, order.CouponID); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "could not make the order", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "Successfully made the order", nil, nil)
	c.JSON(http.StatusOK, successRes)
}

func (i *OrderHandler) GetOrders(c *gin.Context) {

	idString := c.Query("order_id")
	order_id, err := strconv.Atoi(idString)

	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "check your id again", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	orders, err := i.orderUseCase.GetOrders(order_id)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "could not retrieve orders", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "Successfully retrieved all orders", orders, nil)
	c.JSON(http.StatusOK, successRes)

}

func (i OrderHandler) CancelOrder(c *gin.Context) {
	idString := c.Query("order_id")
	orderID, err := strconv.Atoi(idString)

	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "check your id again", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	err = i.orderUseCase.CancelOrder(orderID)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "could not cancel the order", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Order successfully canceled", nil, nil)
	c.JSON(http.StatusOK, successRes)
}

func (i *OrderHandler) GetAllOrders(c *gin.Context) {

	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "page number not in correct format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("count", "10"))
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "page count not in right format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}

	id, _ := c.Get("id")
	UserID, _ := id.(int)

	orders, err := i.orderUseCase.GetAllOrders(UserID, page, pageSize)

	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "could not retrieve orders", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "Successfully retrieved all orders", orders, nil)
	c.JSON(http.StatusOK, successRes)

}

func (i *OrderHandler) GetAdminOrders(c *gin.Context) {

	pageStr := c.Query("page")
	page, err := strconv.Atoi(pageStr)

	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "page number not in right format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	orders, err := i.orderUseCase.GetAdminOrders(page)

	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "could not retrieve orders", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "Successfully retrieved all orders", orders, nil)
	c.JSON(http.StatusOK, successRes)

}

func (i *OrderHandler) ApproveOrder(c *gin.Context) {
	id := c.Query("order_id")
	orderID, err := strconv.Atoi(id)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Invalid order ID format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	err = i.orderUseCase.OrdersStatus(orderID)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not approve order", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully approved order", nil, nil)
	c.JSON(http.StatusOK, successRes)
}

func (o *OrderHandler) ReturnOrder(c *gin.Context) {
	id := c.Query("order_id")
	orderID, err := strconv.Atoi(id)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "Invalid order ID format", nil, err)
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	fmt.Println("order id at handler", orderID)

	err = o.orderUseCase.ReturnOrder(orderID)
	if err != nil {
		errRes := response.ClientResponse(http.StatusInternalServerError, "Failed to process order return", nil, err)
		c.JSON(http.StatusInternalServerError, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Order successfully returned", nil, nil)
	c.JSON(http.StatusOK, successRes)
}
