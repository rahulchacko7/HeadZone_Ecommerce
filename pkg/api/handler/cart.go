package handler

import (
	"HeadZone/pkg/usecase/interfaces"
	"HeadZone/pkg/utils/models"

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

func (i *CartHandler) AddToCart(c *gin.Context) {

	var model models.AddToCart
	if err := c.BindJSON(&model); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	if err := i.usecase.AddToCart(model.UserID, model.InventoryID); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not add the Cart", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully added To cart", nil, nil)
	c.JSON(http.StatusOK, successRes)

}
