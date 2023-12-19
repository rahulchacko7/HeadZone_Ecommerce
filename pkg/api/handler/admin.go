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
