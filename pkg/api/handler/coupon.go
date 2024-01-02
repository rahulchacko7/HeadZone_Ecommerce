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

// AddCoupon adds a new coupon.
// @Summary Add a new coupon
// @Description Adds a new coupon with provided details
// @Tags Admin Coupon Management
// @Accept json
// @Produce json
// @security BearerTokenAuth
// @Param CouponDetails body models.CouponResponse true "Coupon details in JSON format"
// @Success 200 {object} response.Response "Successfully updated the coupon"
// @Failure 400 {object} response.Response "Invalid request or incorrect format"
// @Router /admin/coupon [post]
func (handler *CouponHandler) AddCoupon(c *gin.Context) {
	var cp models.CouponResponse

	if err := c.BindJSON(&cp); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "Fields are provided in the wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	coupon, err := handler.CouponUseCase.AddCoupon(cp.CouponName, cp.Status, cp.DiscountRate)

	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "Error updating coupon", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully updated coupon", coupon, nil)
	c.JSON(http.StatusOK, successRes)
}

// GetCoupons retrieves all available coupons.
// @Summary Retrieve all coupons
// @Description Retrieves all available coupons
// @Tags Admin Coupon Management
// @Accept json
// @Produce json
// @security BearerTokenAuth
// @Success 200 {object} response.Response "Successfully retrieved all coupons"
// @Failure 400 {object} response.Response "Error retrieving coupons"
// @Router /admin/coupon [get]
func (handler *CouponHandler) GetCoupons(c *gin.Context) {
	coupons, err := handler.CouponUseCase.GetCoupon()
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "fields are provided in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
	}
	successRes := response.ClientResponse(http.StatusOK, "sucessfully retrived all records", coupons, nil)
	c.JSON(http.StatusOK, successRes)
}

// UpdateCoupon updates an existing coupon by ID.
// @Summary Update an existing coupon by ID
// @Description Update an existing coupon's details such as name, status, and discount rate
// @Tags Admin Coupon Management
// @Accept json
// @Produce json
// @security BearerTokenAuth
// @Param id query integer true "Coupon ID to update"
// @Param CouponDetails body models.CouponResponse true "Coupon details to update"
// @Success 200 {object} models.CouponResponse "Successfully updated coupon"
// @Failure 400 {object} response.Response "Error updating coupon"
// @Router /admin/coupon [patch]
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

	coupon, err := handler.CouponUseCase.UpdateCoupon(CId, cp.CouponName, cp.Status, cp.DiscountRate)

	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "Error updating coupon", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully updated coupon", coupon, nil)
	c.JSON(http.StatusOK, successRes)
}

// GetAllCoupons retrieves all coupons.
// @Summary Retrieve all coupons
// @Description Fetches all available coupons
// @Tags Admin Coupon Management
// @Accept json
// @Produce json
// @security BearerTokenAuth
// @Success 200 {object} models.CouponResponse "Success"
// @Failure 400 {object} response.Response "Error retrieving coupons"
// @Router /admin/coupon [get]
func (handler *CouponHandler) GetAllCoupons(c *gin.Context) {
	coupons, err := handler.CouponUseCase.GetCoupon()
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "fields are provided in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
	}
	successRes := response.ClientResponse(http.StatusOK, "sucessfully retrived all records", coupons, nil)
	c.JSON(http.StatusOK, successRes)
}
