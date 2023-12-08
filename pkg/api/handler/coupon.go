package handler

import (
	"HeadZone/pkg/usecase/interfaces"
	"HeadZone/pkg/utils/models"
	"HeadZone/pkg/utils/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CouponHandler struct {
	CouponUseCase interfaces.CouponUseCase
}

func NewCouponHandler(usecase interfaces.CouponUseCase) *CouponHandler {
	return &CouponHandler{
		CouponUseCase: usecase,
	}
}

func (handler *CouponHandler) AddCoupon(c *gin.Context) {
	var cp models.CouponResponse

	if err := c.BindJSON(&cp); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "Fields are provided in the wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	coupon, err := handler.CouponUseCase.AddCoupon(cp.CouponName, cp.Status, cp.DiscountPercentage, cp.MinimumPrice)

	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "Error updating coupon", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully updated coupon", coupon, nil)
	c.JSON(http.StatusOK, successRes)
}

func (handler *CouponHandler) GetCoupons(c *gin.Context) {
	coupons, err := handler.CouponUseCase.GetCoupon()
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "fields are provided in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
	}
	successRes := response.ClientResponse(http.StatusOK, "sucessfully retrived all records", coupons, nil)
	c.JSON(http.StatusOK, successRes)
}
func (handler *CouponHandler) UpdateCoupon(c *gin.Context) {
	CouponId := c.Query("id")
	CId, err := strconv.Atoi(CouponId)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "Invalid ID format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	var cp models.CouponResponse

	if err := c.BindJSON(&cp); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "Fields are provided in the wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	coupon, err := handler.CouponUseCase.UpdateCoupon(CId, cp.CouponName, cp.Status, cp.DiscountPercentage, cp.MinimumPrice)

	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "Error updating coupon", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully updated coupon", coupon, nil)
	c.JSON(http.StatusOK, successRes)
}

func (cn *CouponHandler) RedeemCoupon(c *gin.Context) {
	userID, ok := c.Get("id")
	if !ok {
		errorRes := response.ClientResponse(http.StatusBadRequest, "User ID not found in context", nil, "User ID not found")
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	var CouponName models.CouponDetails
	if err := c.ShouldBindJSON(&CouponName); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not bind the coupon", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	userIDInt, ok := userID.(int)
	if !ok {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Invalid user ID type in context", nil, "Invalid user ID type")
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	err := cn.CouponUseCase.RedeemCoupon(CouponName.CouponName, userIDInt)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusInternalServerError, "Coupon could not be added", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusCreated, "Coupon added successfully", nil, nil)
	c.JSON(http.StatusCreated, successRes)
}
