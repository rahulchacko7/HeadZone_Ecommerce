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
	CouponName := c.Query("coupon_name")
	CouponStatus, err := strconv.ParseBool(c.Query("coupon_status"))
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "Invalid coupon status format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	Discount, err := strconv.Atoi(c.Query("discount"))
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "Invalid discount format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	MinPurchase, err := strconv.ParseFloat(c.Query("min_purchase"), 64)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "Invalid min purchase format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	coupon, err := handler.CouponUseCase.AddCoupon(CouponName, CouponStatus, Discount, MinPurchase)

	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "Error adding coupon", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully Added", coupon, nil)
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
