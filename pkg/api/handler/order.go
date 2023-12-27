package handler

import (
	"HeadZone/pkg/usecase/interfaces"
	models "HeadZone/pkg/utils/models"
	response "HeadZone/pkg/utils/response"
	"errors"
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

// OrderItemsFromCart creates an order from the items in the user's cart.
// @Summary Create an order from cart items
// @Description Allows a user to create an order from items in their cart
// @Tags User Order Management
// @Accept json
// @Produce json
// @security BearerTokenAuth
// @Param id header integer true "User ID"
// @Param order body models.Order true "Order details"
// @Success 200 {object} response.Response "Success"
// @Failure 400 {object} response.Response "Invalid request format"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /user/check-out [post]
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

// GetOrders retrieves orders based on the provided order ID.
// @Summary Retrieve orders
// @Description Retrieves orders based on the provided order ID
// @Tags User Order Management
// @Accept json
// @Produce json
// @security BearerTokenAuth
// @Param order_id query integer true "Order ID"
// @Success 200 {object} response.Response "Success"
// @Failure 400 {object} response.Response "Invalid request format or missing ID"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /user/profile/orders [get]
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

// CancelOrder cancels an order based on the provided order ID.
// @Summary Cancel an order
// @Description Cancels an order based on the provided order ID
// @Tags Admin Order Management
// @Accept json
// @Produce json
// @security BearerTokenAuth
// @Param order_id query integer true "Order ID"
// @Success 200 {object} response.Response "Success"
// @Failure 400 {object} response.Response "Invalid request format or missing ID"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /user/profile/orders [put]
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

// GetAllOrders retrieves paginated orders for a specific user.
// @Summary Get all orders
// @Description Retrieves paginated orders for a specific user
// @Tags Admin Order Management
// @Accept json
// @Produce json
// @security BearerTokenAuth
// @Param page query integer false "Page number (default: 1)"
// @Param count query integer false "Number of items per page (default: 10)"
// @Success 200 {object} response.Response "Success"
// @Failure 400 {object} response.Response "Invalid request format or missing ID"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /user/profile/orders/all [get]
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

// GetAdminOrders retrieves paginated orders for admin view.
// @Summary Get admin orders
// @Description Retrieves paginated orders for admin view
// @Tags Admin Order Management
// @Accept json
// @Produce json
// @security BearerTokenAuth
// @Param page query integer true "Page number"
// @Success 200 {object} response.Response "Success"
// @Failure 400 {object} response.Response "Invalid request format"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /admin/orders [get]
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

// ApproveOrder approves an order by its ID.
// @Summary Approve an order
// @Description Approves an order by its ID
// @Tags Admin Order Management
// @Accept json
// @Produce json
// @security BearerTokenAuth
// @Param order_id query integer true "Order ID"
// @Success 200 {object} response.Response "Success"
// @Failure 400 {object} response.Response "Invalid request format"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /admin/orders/status [GET]
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

// ReturnOrder handles the return of an order by its ID.
// @Summary Return an order
// @Description Returns an order by its ID
// @Tags Orders
// @Accept json
// @Produce json
// @Param order_id query integer true "Order ID"
// @Success 200 {object} response.Response "Success"
// @Failure 400 {object} response.Response "Invalid request format"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /orders/return [put]
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

// PrintInvoice generates and provides a PDF invoice for a specific order ID.
// @Summary Print an invoice
// @Description Generates and provides a PDF invoice for a specific order ID
// @Tags User Invoice
// @Accept json
// @Produce pdf
// @security BearerTokenAuth
// @Param order_id query integer true "Order ID"
// @Success 200 {string} pdf "Invoice PDF file"
// @Failure 400 {object} response.Response "Invalid request format"
// @Failure 502 {object} response.Response "Bad Gateway error"
// @Router /user/check-out/print [get]
func (O *OrderHandler) PrintInvoice(c *gin.Context) {
	orderId := c.Query("order_id")
	orderIdInt, err := strconv.Atoi(orderId)
	if err != nil {
		err = errors.New("error in coverting order id" + err.Error())
		errRes := response.ClientResponse(http.StatusBadGateway, "error in reading the order id", nil, err)
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	pdf, err := O.orderUseCase.PrintInvoice(orderIdInt)
	fmt.Println("error ", err)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadGateway, "error in printing the invoice", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	c.Header("Content-Disposition", "attachment;filename=invoice.pdf")

	pdfFilePath := "salesReport/invoice.pdf"

	err = pdf.OutputFileAndClose(pdfFilePath)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadGateway, "error in printing invoice", nil, err)
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	c.Header("Content-Disposition", "attachment; filename=sales_report.pdf")
	c.Header("Content-Type", "application/pdf")

	c.File(pdfFilePath)

	c.Header("Content-Type", "application/pdf")

	err = pdf.Output(c.Writer)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadGateway, "error in printing invoice", nil, err)
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "the request was succesful", pdf, nil)
	c.JSON(http.StatusOK, successRes)
}
