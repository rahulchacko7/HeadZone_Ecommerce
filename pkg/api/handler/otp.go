package handler

import (
	"HeadZone/pkg/usecase/interfaces"
	models "HeadZone/pkg/utils/models"
	"HeadZone/pkg/utils/response"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OtpHandler struct {
	otpUseCase interfaces.OtpUseCase
}

func NewOtpHandler(useCase interfaces.OtpUseCase) *OtpHandler {
	return &OtpHandler{
		otpUseCase: useCase,
	}
}

// SendOTP sends an OTP (One-Time Password) to a specified phone number.
// @Summary Send OTP
// @Description Sends an OTP (One-Time Password) to the provided phone number
// @Tags User OTP
// @Accept json
// @Produce json
// @Param phone body models.OTPData true "Phone number to receive OTP"
// @Success 200 {object} response.Response "OTP sent successfully"
// @Failure 400 {object} response.Response "Invalid request format or phone number"
// @Failure 502 {object} response.Response "Bad Gateway error"
// @Router /user/otplogin [post]
func (ot *OtpHandler) SendOTP(c *gin.Context) {

	var phone models.OTPData
	if err := c.BindJSON(&phone); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
	}

	err := ot.otpUseCase.SendOTP(phone.PhoneNumber)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not send OTP", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "OTP sent successfully", nil, nil)
	c.JSON(http.StatusOK, successRes)

}

// VerifyOTP verifies the provided OTP (One-Time Password).
// @Summary Verify OTP
// @Description Verifies the provided OTP (One-Time Password)
// @Tags User OTP
// @Accept json
// @Produce json
// @Param code body models.VerifyData true "Data to verify OTP"
// @Success 200 {object} response.Response "OTP verified successfully"
// @Failure 400 {object} response.Response "Invalid request format or OTP"
// @Failure 502 {object} response.Response "Bad Gateway error"
// @Router /user/verifyotp [post]
func (ot *OtpHandler) VerifyOTP(c *gin.Context) {
	fmt.Println(1)

	var code models.VerifyData
	if err := c.BindJSON(&code); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	fmt.Println(2)
	users, err := ot.otpUseCase.VerifyOTP(code)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not verify OTP", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully verified OTP", users, nil)
	c.JSON(http.StatusOK, successRes)

}
