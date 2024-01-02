package handler

import (
	"HeadZone/pkg/usecase/interfaces"
	"strconv"

	"HeadZone/pkg/utils/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	usecase interfaces.CartUseCase
}

func NewCartHandler(usecase interfaces.CartUseCase) *CartHandler {
	return &CartHandler{
		usecase: usecase,
	}
}

// AddToCart adds a product to the user's cart.
// @Summary Add a product to the cart
// @Description Adds a selected product to the user's cart based on provided parameters
// @Tags User Cart Management
// @Accept json
// @Produce json
// @security BearerTokenAuth
// @Param inventory_id query integer true "Inventory ID of the product to be added"
// @Param quantity query integer true "Quantity of the product to be added"
// @Security ApiKeyAuth
// @Success 200 {object} interface{} "Successfully added to cart"
// @Failure 400 {object} response.Response "Invalid request or incorrect format"
// @Failure 401 {object} response.Response "Unauthorized"
// @Router /user/cart [post]
func (i *CartHandler) AddToCart(c *gin.Context) {
	idString, _ := c.Get("id")
	UserID, _ := idString.(int)

	InventoryID, err := strconv.Atoi(c.Query("inventory_id"))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Inventory Id not in right format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	qty, err := strconv.Atoi(c.Query("quantity"))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "check parameters properly", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	if err := i.usecase.AddToCart(UserID, InventoryID, qty); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not add to the cart", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully added to cart", nil, nil)
	c.JSON(http.StatusOK, successRes)
}

// CheckOut fetches the products in the user's cart for checkout.
// @Summary Get products for checkout
// @Description Retrieves the products in the user's cart for checkout
// @Tags User Cart Management
// @Accept json
// @Produce json
// @security BearerTokenAuth
// @Success 200 {object} interface{} "Successfully got all records"
// @Failure 400 {object} response.Response "Invalid request or incorrect format"
// @Failure 401 {object} response.Response "Unauthorized"
// @Router /user/check-out [get]
func (i *CartHandler) CheckOut(c *gin.Context) {
	idString, _ := c.Get("id")
	id, _ := idString.(int)

	products, err := i.usecase.CheckOut(id)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "could not open checkout", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "Successfully got all records", products, nil)
	c.JSON(http.StatusOK, successRes)
}
