package handler

import (
	"HeadZone/pkg/helper"
	services "HeadZone/pkg/usecase/interfaces"
	models "HeadZone/pkg/utils/models"
	response "HeadZone/pkg/utils/response"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt"
)

type AdminHandler struct {
	adminUseCase services.AdminUseCase
}

func NewAdminHandler(usecase services.AdminUseCase) *AdminHandler {
	return &AdminHandler{
		adminUseCase: usecase,
	}
}

// @Summary Authenticate admin
// @Description Authenticate an admin user
// @Tags Admin
// @Accept json
// @Produce json
// @Param request body models.AdminLogin true "Admin login details in JSON format"
// @Success 200 {object} domain.TokenAdmin "Admin authenticated successfully"
// @Failure 400 {object} response.Response "Invalid request or incorrect format"
// @Failure 400 {object} response.Response "Constraints not satisfied"
// @Failure 400 {object} response.Response "Unable to authenticate user"
// @Router /admin [post]
func (ad *AdminHandler) LoginHandler(c *gin.Context) {

	var adminDetails models.AdminLogin
	if err := c.BindJSON(&adminDetails); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "details not in correct format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	err := validator.New().Struct(adminDetails)

	if err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, "constraints not satisfied", nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	admin, err := ad.adminUseCase.LoginHandler(adminDetails)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "cannot authenticate user", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	c.Set("Access", admin.AccessToken)
	// c.Set("Refresh", admin.RefreshToken)

	successRes := response.ClientResponse(http.StatusOK, "Admin authenticated successfully", admin, nil)
	c.JSON(http.StatusOK, successRes)

}

// BlockUser handles the blocking of a user by ID.
// @Summary Block a user
// @Description Block a user based on their ID
// @Tags Admin User Management
// @Accept json
// @Produce json
// @security BearerTokenAuth
// @Param id query string true "User ID to block"
// @Success 200 {object} response.Response "Successfully blocked the user"
// @Failure 400 {object} response.Response "User could not be blocked"
// @Router /admin/users/block [post]
func (ad *AdminHandler) BlockUser(c *gin.Context) {

	id := c.Query("id")
	err := ad.adminUseCase.BlockUser(id)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "user could not be blocked", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully blocked the user", nil, nil)
	c.JSON(http.StatusOK, successRes)

}

// UnBlockUser handles the unblocking of a user by their ID.
// @Summary Unblock a user
// @Description Unblock a user based on their ID
// @Tags Admin User Management
// @Accept json
// @Produce json
// @security BearerTokenAuth
// @Param id query string true "User ID to be unblocked"
// @Success 200 {object} response.Response "Successfully unblocked the user"
// @Failure 400 {object} response.Response "User could not be unblocked"
// @Router /admin/users/unblock [post]
func (ad *AdminHandler) UnBlockUser(c *gin.Context) {

	id := c.Query("id")
	err := ad.adminUseCase.UnBlockUser(id)

	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "user could not be unblocked", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully unblocked the user", nil, nil)
	c.JSON(http.StatusOK, successRes)
}

// GetUsers retrieves a list of users based on the specified page number.
// @Summary Retrieve users
// @Description Get a list of users based on the provided page number
// @Tags Admin User Management
// @Accept json
// @Produce json
// @security BearerTokenAuth
// @Param page query integer true "Page number for user retrieval"
// @Success 200 {object} response.Response "Successfully retrieved the users"
// @Failure 400 {object} response.Response "Page number not in the right format or could not retrieve records"
// @Router /admin/users [get]
func (ad *AdminHandler) GetUsers(c *gin.Context) {

	pageStr := c.Query("page")
	page, err := strconv.Atoi(pageStr)

	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "page number not in right format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	users, err := ad.adminUseCase.GetUsers(page)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "could not retrieve records", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "Successfully retrieved the users", users, nil)
	c.JSON(http.StatusOK, successRes)

}

func (a *AdminHandler) ValidateRefreshTokenAndCreateNewAccess(c *gin.Context) {

	refreshToken := c.Request.Header.Get("RefreshToken")

	// Check if the refresh token is valid.
	_, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		return []byte("refreshsecret"), nil
	})
	if err != nil {
		// The refresh token is invalid.
		c.AbortWithError(401, errors.New("refresh token is invalid:user have to login again"))
		return
	}

	claims := &helper.AuthCustomClaims{
		Role: "admin",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	newAccessToken, err := token.SignedString([]byte("accesssecret"))
	if err != nil {
		c.AbortWithError(500, errors.New("error in creating new access token"))
	}

	c.JSON(200, newAccessToken)
}

// NewPaymentMethod creates a new payment method.
//
// @Summary Create a new payment method
// @Description Add a new payment method to the system
// @Tags Admin Payment Method
// @Accept json
// @Produce json
// @security BearerTokenAuth
// @Param request body models.NewPaymentMethod true "New payment method details in JSON format"
// @Success 200 {string} string "Successfully added Payment Method"
// @Failure 400 {object} response.Response "Invalid request or incorrect format"
// @Router /admin/payment-method [post]
func (i *AdminHandler) NewPaymentMethod(c *gin.Context) {

	var method models.NewPaymentMethod
	if err := c.BindJSON(&method); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	err := i.adminUseCase.NewPaymentMethod(method.PaymentMethod)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not add the payment method", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully added Payment Method", nil, nil)
	c.JSON(http.StatusOK, successRes)

}

// ListPaymentMethods retrieves a list of available payment methods.
//
// @Summary Retrieve available payment methods
// @Description Get a list of all available payment methods in the system
// @Tags Admin Payment Method
// @Produce json
// @security BearerTokenAuth
// @Success 200 {object} []models.PaymentMethod "List of payment methods"
// @Failure 400 {object} response.Response "Invalid request or incorrect format"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /admin/payment-methods [get]
func (a *AdminHandler) ListPaymentMethods(c *gin.Context) {

	categories, err := a.adminUseCase.ListPaymentMethods()
	if err != nil {
		errorRes := response.ClientResponse(http.StatusInternalServerError, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully got all payment methods", categories, nil)
	c.JSON(http.StatusOK, successRes)

}

// DeletePaymentMethod deletes a payment method by ID.
//
// @Summary Delete a payment method
// @Description Delete a payment method by its ID
// @Tags Admin Payment Method
// @security BearerTokenAuth
// @Param id query integer true "Payment method ID to delete"
// @Success 200 {object} response.Response "Successfully deleted the payment method"
// @Failure 400 {object} response.Response "Invalid request or incorrect format"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /admin/payment-methods [delete]
func (a *AdminHandler) DeletePaymentMethod(c *gin.Context) {

	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	err = a.adminUseCase.DeletePaymentMethod(id)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusInternalServerError, "error in deleting data", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully deleted the Category", nil, nil)
	c.JSON(http.StatusOK, successRes)

}

// DashBoard retrieves the dashboard details.
//
// @Summary Retrieve dashboard details
// @Description Get details for the dashboard
// @Tags Admin Dashboard
// @security BearerTokenAuth
// @Success 200 {object} response.Response "Successfully received the dashboard details"
// @Failure 400 {object} response.Response "Invalid request or incorrect format"
// @Router /admin/dashboard [get]
func (a *AdminHandler) DashBoard(c *gin.Context) {
	dashBoard, err := a.adminUseCase.DashBoard()
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "error in getting dashboard details", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	sucessRes := response.ClientResponse(http.StatusOK, "succesfully recevied all records", dashBoard, nil)
	c.JSON(http.StatusOK, sucessRes)
}

// SalesByDate retrieves sales details by date and generates reports in PDF or Excel format.
//
// @Summary Retrieve sales details by date
// @Description Get sales details based on the provided year, month, and day parameters. Generate reports in PDF or Excel format.
// @Tags Admin Dashboard
// @Accept json
// @Produce json
// @security BearerTokenAuth
// @Param year query int true "Year (YYYY)"
// @Param month query int true "Month (1-12)"
// @Param day query int true "Day (1-31)"
// @Param download query string true "Specify 'pdf' or 'excel' to download the report in PDF or Excel format respectively"
// @Success 200 {object} response.Response "Successfully received the sales details"
// @Failure 400 {object} response.Response "Invalid request or incorrect format"
// @Router /admin/salesbydate [get]
func (a *AdminHandler) SalesByDate(c *gin.Context) {
	year := c.Query("year")
	yearInt, err := strconv.Atoi(year)

	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "error in getting year", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	month := c.Query("month")
	monthInt, err := strconv.Atoi(month)

	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "error in getting month", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	day := c.Query("day")
	dayInt, err := strconv.Atoi(day)

	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "error in getting day", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	body, err := a.adminUseCase.SalesByDate(dayInt, monthInt, yearInt)

	fmt.Println("body handler", dayInt)
	fmt.Println("body handler", monthInt)
	fmt.Println("body handler", yearInt)

	fmt.Println("body ", body)

	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "error in getting sales details", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	download := c.Query("download")
	if download == "pdf" {
		pdf, err := a.adminUseCase.PrintSalesReport(body)
		if err != nil {
			errRes := response.ClientResponse(http.StatusBadGateway, "error in printing sales report", nil, err)
			c.JSON(http.StatusBadRequest, errRes)
			return
		}
		c.Header("Content-Disposition", "attachment;filename=totalsalesreport.pdf")

		pdfFilePath := "salesReport/totalsalesreport.pdf"

		err = pdf.OutputFileAndClose(pdfFilePath)
		if err != nil {
			errRes := response.ClientResponse(http.StatusBadGateway, "error in printing sales report", nil, err)
			c.JSON(http.StatusBadRequest, errRes)
			return
		}

		c.Header("Content-Disposition", "attachment; filename=total_sales_report.pdf")
		c.Header("Content-Type", "application/pdf")

		c.File(pdfFilePath)

		c.Header("Content-Type", "application/pdf")

		err = pdf.Output(c.Writer)
		if err != nil {
			errRes := response.ClientResponse(http.StatusBadGateway, "error in printing sales report", nil, err)
			c.JSON(http.StatusBadRequest, errRes)
			return
		}
	} else {
		fmt.Println("body ", body)
		excel, err := helper.ConvertToExel(body)
		if err != nil {
			errRes := response.ClientResponse(http.StatusBadGateway, "error in printing sales report", nil, err)
			c.JSON(http.StatusBadRequest, errRes)
			return
		}

		fileName := "sales_report.xlsx"

		c.Header("Content-Disposition", "attachment; filename="+fileName)
		c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

		if err := excel.Write(c.Writer); err != nil {
			errRes := response.ClientResponse(http.StatusBadGateway, "Error in serving the sales report", nil, err)
			c.JSON(http.StatusBadRequest, errRes)
			return
		}
	}

	succesRes := response.ClientResponse(http.StatusOK, "success", body, nil)
	c.JSON(http.StatusOK, succesRes)
}

// CustomSalesReport generates a custom sales report based on the provided start and end dates.
//
// @Summary Generate custom sales report
// @Description Generates a sales report within the specified date range. Requires 'start' and 'end' dates in the format 'DD-MM-YYYY'.
// @Tags Admin Dashboard
// @Accept json
// @Produce json
// @security BearerTokenAuth
// @Param start query string true "Start date (DD-MM-YYYY)"
// @Param end query string true "End date (DD-MM-YYYY)"
// @Success 200 {object} response.Response "Custom sales report retrieved successfully"
// @Failure 400 {object} response.Response "Invalid request or incorrect format"
// @Router /admin/customreport [get]
func (ad *AdminHandler) CustomSalesReport(c *gin.Context) {
	startDateStr := c.Query("start")
	endDateStr := c.Query("end")
	if startDateStr == "" || endDateStr == "" {
		err := response.ClientResponse(http.StatusBadRequest, "start or end date is empty", nil, "Empty date string")
		c.JSON(http.StatusBadRequest, err)
		return
	}
	startDate, err := time.Parse("02-01-2006", startDateStr)
	if err != nil {
		err := response.ClientResponse(http.StatusBadRequest, "start date conversion failed", nil, err.Error())
		c.JSON(http.StatusBadRequest, err)
		return
	}
	endDate, err := time.Parse("02-01-2006", endDateStr)
	if err != nil {
		err := response.ClientResponse(http.StatusBadRequest, "end date conversion failed", nil, err.Error())
		c.JSON(http.StatusBadRequest, err)
		return
	}

	if startDate.After(endDate) {
		err := response.ClientResponse(http.StatusBadRequest, "start date is after end date", nil, "Invalid date range")
		c.JSON(http.StatusBadRequest, err)
		return
	}

	report, err := ad.adminUseCase.CustomSalesReportByDate(startDate, endDate)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusInternalServerError, "sales report could not be retrieved", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errorRes)
		return
	}

	success := response.ClientResponse(http.StatusOK, "custom report retrieved successfully", report, nil)
	c.JSON(http.StatusOK, success)
}
