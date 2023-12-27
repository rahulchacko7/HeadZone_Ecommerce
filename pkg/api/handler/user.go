package handler

import (
	services "HeadZone/pkg/usecase/interfaces"
	"HeadZone/pkg/utils/models"
	"HeadZone/pkg/utils/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type UserHandler struct {
	userUseCase services.UserUseCase
}

type Response struct {
	ID      uint   `copier:"must"`
	Name    string `copier:"must"`
	Surname string `copier:"must"`
}

func NewUserHandler(usecase services.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: usecase,
	}
}

// UserSignUp handles user sign-up functionality.
// @Summary Register a new user
// @Description Create a new user account
// @Tags users
// @Accept json
// @Produce json
// @Param request body models.UserDetails true "User details in JSON format"
// @Success 201 {object} models.UserDetails "User signed up successfully"
// @Failure 400 {object} models.TokenUsers "Invalid request or constraints not satisfied"
// @Failure 500 {object}  models.TokenUsers "Internal server error"
// @Router /user/signup [post]
func (u *UserHandler) UserSignUp(c *gin.Context) {

	var user models.UserDetails
	if err := c.BindJSON(&user); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	err := validator.New().Struct(user)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "constraints not satisfied", nil, err.Error())
		c.JSON(http.StatusBadRequest,
			errRes)
		return
	}

	userCreated, err := u.userUseCase.UserSignUp(user)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "User could not signed up", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusCreated, "User successfully signed up", userCreated, nil)
	c.JSON(http.StatusCreated, successRes)

}

// LoginHandler handles the user login functionality.
// @Summary Log in a user
// @Description Logs in a user with provided credentials
// @Tags users
// @Accept json
// @Produce json
// @Param request body models.UserLogin true "User login details in JSON format"
// @Success 200 {object} models.UserDetails "User details logged in successfully"
// @Failure 400 {object} response.Response "Invalid request or constraints not satisfied"
// @Failure 401 {object} models.UserDetails "Unauthorized: Invalid credentials"
// @Router /user/login [post]
func (u *UserHandler) LoginHandler(c *gin.Context) {

	var user models.UserLogin

	if err := c.BindJSON(&user); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	err := validator.New().Struct(user)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "constraints not satisfied", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	user_details, err := u.userUseCase.LoginHandler(user)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "User could not be logged in", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "User successfully logged in", user_details, nil)
	c.JSON(http.StatusOK, successRes)

}

// GetUserDetails handles the retrieval of user details.
// @Summary Get user details
// @Description Retrieve user details by ID
// @Tags users
// @Accept json
// @Produce json
// @security BearerTokenAuth
// @Success 200 {object} models.UserDetailsResponse "User details retrieved successfully"
// @Failure 400 {object} response.Response "Failed to retrieve user details"
// @Router /user/profile{id} [get]
func (i *UserHandler) GetUserDetails(c *gin.Context) {
	idString, _ := c.Get("id")
	id, _ := idString.(int)

	details, err := i.userUseCase.GetUserDetails(id)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "could not retrieve records", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "Successfully got all records", details, nil)
	c.JSON(http.StatusOK, successRes)
}

// AddAddress handles adding a new address for a user.
// @Summary Add a new address
// @Description Adds a new address for the user
// @Tags users
// @Accept json
// @Produce json
// @security BearerTokenAuth
// @Param request body models.AddAddress true "Address details in JSON format"
// @Success 200 {string} string "Successfully added address"
// @Failure 400 {object} response.Response "Invalid request or address addition failed"
// @Router /user/{id}/address [post]
func (i *UserHandler) AddAddress(c *gin.Context) {
	id, _ := c.Get("id")

	var address models.AddAddress
	if err := c.BindJSON(&address); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	if err := i.userUseCase.AddAddress(id.(int), address); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not add the address", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully added address", nil, nil)
	c.JSON(http.StatusOK, successRes)

}

// GetAddresses handles the retrieval of addresses for a user.
// @Summary Retrieve addresses for a user
// @Description Get addresses associated with a user ID
// @Tags addresses
// @Accept json
// @Produce json
// @security BearerTokenAuth
// @Success 200 {object} []models.Address "Addresses retrieved successfully"
// @Failure 400 {object} response.Response "Could not retrieve records or invalid request"
// @Router /user/profile/addresses/{id} [get]
func (i *UserHandler) GetAddresses(c *gin.Context) {
	idString, _ := c.Get("id")
	id, _ := idString.(int)

	addresses, err := i.userUseCase.GetAddresses(id)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "could not retrieve records", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "Successfully got all records", addresses, nil)
	c.JSON(http.StatusOK, successRes)
}

// EditDetails handles the editing of user details.
// @Summary Edit user details
// @Description Edit user details based on the provided information
// @Tags users
// @Accept json
// @Produce json
// @security BearerTokenAuth
// @Param request body models.EditDetailsResponse true "User details to be updated"
// @Success 201 {object} models.EditDetailsResponse "Updated user details"
// @Failure 400 {object} response.Response "Invalid request or error updating values"
// @Router /user/profile/{id}/edit [put]
func (i *UserHandler) EditDetails(c *gin.Context) {

	idString, _ := c.Get("id")
	id, _ := idString.(int)

	var model models.EditDetailsResponse
	if err := c.BindJSON(&model); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	body, err := i.userUseCase.EditDetails(id, model)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "error updating the values", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusCreated, "addresses fetched succesfully", body, nil)

	c.JSON(http.StatusCreated, successRes)
}

// ChangePassword allows a user to change their password.
// @Summary Change Password
// @Description Allows a user to update their password
// @Tags Users
// @Accept json
// @Produce json
// @security BearerTokenAuth
// @Param id header integer true "User ID"
// @Param ChangePassword body models.ChangePassword true "Change Password Request"
// @Success 200 {object} response.ClientResponse "Password changed successfully"
// @Failure 400 {object} response.ClientResponse "Invalid request format or password change failure"
// @Router /user/change-password [post]
func (i *UserHandler) ChangePassword(c *gin.Context) {

	idString, _ := c.Get("id")
	id, _ := idString.(int)

	var ChangePassword models.ChangePassword
	if err := c.BindJSON(&ChangePassword); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	if err := i.userUseCase.ChangePassword(id, ChangePassword.Oldpassword, ChangePassword.Password, ChangePassword.Repassword); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not change the password", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "password changed Successfully ", nil, nil)
	c.JSON(http.StatusOK, successRes)

}

// GetCart retrieves the user's cart contents.
// @Summary Get User Cart
// @Description Retrieves the products in the user's cart
// @Tags User Cart Management
// @Accept json
// @Produce json
// @security BearerTokenAuth
// @Param id header integer true "User ID"
// @Success 200 {object} response.Response "User's cart retrieved successfully"
// @Failure 400 {object} response.Response "Failed to retrieve the cart"
// @Router /user/cart [get]
func (i *UserHandler) GetCart(c *gin.Context) {
	idString, _ := c.Get("id")
	id, _ := idString.(int)

	products, err := i.userUseCase.GetCart(id)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "could not retrieve cart", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "Successfully got all products in cart", products, nil)
	c.JSON(http.StatusOK, successRes)
}

// RemoveFromCart removes a product from the user's cart.
// @Summary Remove Product from Cart
// @Description Removes a specific product from the user's cart by cart ID and inventory ID
// @Tags User Cart Managament
// @Accept json
// @Produce json
// @security BearerTokenAuth
// @Param cart_id query integer true "Cart ID"
// @Param inventory_id query integer true "Inventory ID"
// @Success 200 {object} response.ClientResponse "Product removed successfully from the cart"
// @Failure 400 {object} response.ClientResponse "Failed to remove the product from the cart"
// @Router /user/cart [delete]
func (i *UserHandler) RemoveFromCart(c *gin.Context) {

	cartID, err := strconv.Atoi(c.Query("cart_id"))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "check parameters properly", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	InventoryID, err := strconv.Atoi(c.Query("inventory_id"))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "check parameters properly", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	if err := i.userUseCase.RemoveFromCart(cartID, InventoryID); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "could not remove from cart", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully Removed product from cart", nil, nil)
	c.JSON(http.StatusOK, successRes)
}

// UpdateQuantity updates the quantity of a product in the user's cart.
// @Summary Update Product Quantity in Cart
// @Description Updates the quantity of a specific product in the user's cart by ID and inventory ID
// @Tags User Cart Managament
// @Accept json
// @Produce json
// @security BearerTokenAuth
// @Param id query integer true "Product ID"
// @Param inventory query integer true "Inventory ID"
// @Param quantity query integer true "New Quantity"
// @Success 200 {object} response.Response "Quantity updated successfully in the cart"
// @Failure 400 {object} response.Response "Failed to update the quantity in the cart"
// @Router /user/cart [put]
func (i *UserHandler) UpdateQuantity(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "check parameters properly", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	inv, err := strconv.Atoi(c.Query("inventory"))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "check parameters properly", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	qty, err := strconv.Atoi(c.Query("quantity"))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "check parameters properly", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	if err := i.userUseCase.UpdateQuantity(id, inv, qty); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "could not Add the quantity", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully added quantity", nil, nil)
	c.JSON(http.StatusOK, successRes)
}
