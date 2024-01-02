package handler

import (
	"HeadZone/pkg/usecase/interfaces"
	"HeadZone/pkg/utils/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WalletHandler struct {
	WalletUsecase interfaces.WalletUsecase
}

func NewWalletHandler(usecase interfaces.WalletUsecase) *WalletHandler {
	return &WalletHandler{
		WalletUsecase: usecase,
	}
}

// ViewWallet retrieves the wallet details for a specific user.
// @Summary View User's Wallet
// @Description Retrieves the wallet details for a specific user by ID
// @Tags User Wallet Management
// @Accept json
// @Produce json
// @security BearerTokenAuth
// @Param id header integer true "User ID"
// @Success 200 {object} response.Response "Wallet details retrieved successfully"
// @Failure 400 {object} response.Response "Failed to retrieve wallet details"
// @Router /user/wallet [get]
func (handler *WalletHandler) ViewWallet(c *gin.Context) {
	idString, _ := c.Get("id")
	id, _ := idString.(int)

	walletDetails, err := handler.WalletUsecase.GetWallet(id)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "coluld not retrive data", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	sucessRes := response.ClientResponse(http.StatusOK, "Sucessfully retrived wallet", walletDetails, nil)
	c.JSON(http.StatusOK, sucessRes)
}
