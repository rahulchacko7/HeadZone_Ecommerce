package handler

import (
	"HeadZone/pkg/domain"
	"HeadZone/pkg/usecase/interfaces"
	models "HeadZone/pkg/utils/models"
	"fmt"
	"strconv"

	"HeadZone/pkg/utils/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type InventoryHandler struct {
	InventoryUseCase interfaces.InventoryUseCase
}

func NewInventoryHandler(usecase interfaces.InventoryUseCase) *InventoryHandler {
	return &InventoryHandler{
		InventoryUseCase: usecase,
	}
}

// AddInventory adds new inventory.
// @Summary Add new inventory
// @Description Adds new inventory details
// @Tags Admin Inventory Management
// @Accept json
// @Produce json
// @security BearerTokenAuth
// @Param inventory body models.AddInventories true "Inventory object to be added"
// @Success 200 {object} models.Null "Success"
// @Failure 400 {object} response.Response "Error adding inventory"
// @Router /admin/inventory [post]
func (i *InventoryHandler) AddInventory(c *gin.Context) {
	var inventory models.AddInventories

	if err := c.ShouldBindJSON(&inventory); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "form file error", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	fmt.Println("kkkkkkkkkk", inventory)
	_, err := i.InventoryUseCase.AddInventory(inventory)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "Could not add the Inventory", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully added inventory", nil, nil)
	c.JSON(http.StatusOK, successRes)

}

// ListProducts handles the retrieval of products with pagination.
// @Summary List products with pagination
// @Description Retrieves a list of products with pagination support
// @Tags User Product
// @Accept json
// @Produce json
// @Param page query integer false "Page number for pagination (default: 1)"
// @Param per_page query integer false "Number of products per page (default: 5)"
// @Success 200 {array} models.InventoryUserResponse "List of products"
// @Failure 400 {object} response.Response "Invalid request or failed to fetch products"
// @Router /user/products [get]
func (i *InventoryHandler) ListProducts(c *gin.Context) {
	pageNo := c.DefaultQuery("page", "1")
	pageList := c.DefaultQuery("per_page", "5")
	pageNoInt, err := strconv.Atoi(pageNo)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "Invalid page number", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	pageListInt, err := strconv.Atoi(pageList)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "Invalid per_page value", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	productsList, err := i.InventoryUseCase.ListProducts(pageNoInt, pageListInt)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "Failed to fetch products", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Product list", productsList, nil)
	c.JSON(http.StatusOK, successRes)
}

// EditInventory updates an inventory item based on the provided details.
// @Summary Edit an inventory item
// @Description Update an inventory item by ID
// @Tags Admin Inventory Management
// @Accept json
// @Produce json
// @security BearerTokenAuth
// @Param inventory_id query integer true "Inventory ID to update"
// @Param Inventory body domain.Inventory true "Inventory object to update"
// @Success 200 {object} response.Response "Successfully updated inventory"
// @Failure 400 {object} response.Response "Invalid request format or fields in the wrong format"
// @Router /admin/inventory [put]
func (u *InventoryHandler) EditInventory(c *gin.Context) {
	var inventory domain.Inventory

	id := c.Query("inventory_id")
	idInt, err := strconv.Atoi(id)

	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "problems in the id", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	if err := c.BindJSON(&inventory); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "fields are in the wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	modInventory, err := u.InventoryUseCase.EditInventory(inventory, idInt)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "could not edit the product", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "sucessfully edited products", modInventory, nil)
	c.JSON(http.StatusOK, successRes)
}

// DeleteInventory deletes an inventory item by its ID.
// @Summary Delete an inventory item
// @Description Deletes an inventory item by its ID
// @Tags Admin Inventory Management
// @Accept json
// @Produce json
// @security BearerTokenAuth
// @Param id query int true "Inventory ID to delete"
// @Success 200 {object} response.Response "Success"
// @Failure 400 {object} response.Response "Error deleting inventory"
// @Router /admin/inventory [delete]
func (u *InventoryHandler) DeleteInventory(c *gin.Context) {

	inventoryID := c.Query("id")

	err := u.InventoryUseCase.DeleteInventory(inventoryID)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "fields are provided in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "Sucessfully deleted the product", nil, nil)
	c.JSON(http.StatusOK, successRes)
}

// UpdateInventory updates the stock of an inventory item by its ID.
// @Summary Update inventory stock
// @Description Updates the stock of an inventory item by its ID
// @Tags Admin Inventory Management
// @Accept json
// @Produce json
// @security BearerTokenAuth
// @Param Productid body int true "Product ID"
// @Param Stock body int true "New Stock"
// @Success 200 {object} response.Response "Success"
// @Failure 400 {object} response.Response "Error updating inventory stock"
// @Router /admin/inventory/stock [put]
func (i *InventoryHandler) UpdateInventory(c *gin.Context) {

	var p models.InventoryUpdate

	if err := c.BindJSON(&p); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "fileds are provided in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	a, err := i.InventoryUseCase.UpdateInventory(p.Productid, p.Stock)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "Could  not update the inventory stock", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Sucessfully upadated inventory stock", a, nil)
	c.JSON(http.StatusOK, successRes)
}

func (i *InventoryHandler) ShowIndividualProducts(c *gin.Context) {

	id := c.Query("id")
	product, err := i.InventoryUseCase.ShowIndividualProducts(id)

	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "path variables in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Product details retrieved successfully", product, nil)
	c.JSON(http.StatusOK, successRes)
}

// UpdateInventory updates the stock of an inventory item by its ID.
// @Summary Update inventory stock
// @Description Search for a product
// @Tags User Product
// @Accept json
// @Produce json
// @Param Productid body int true "Product ID"
// @Param Stock body int true "New Stock"
// @Success 200 {object} response.Response "Success"
// @Failure 400 {object} response.Response "Error updating inventory stock"
// @Router /user/products/search [put]
func (i *InventoryHandler) SearchProducts(c *gin.Context) {

	var prefix models.SearchItems

	if err := c.ShouldBindJSON(&prefix); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "fields are provided in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	productDetails, err := i.InventoryUseCase.SearchProductsOnPrefix(prefix.ProductName)

	if err != nil {
		errRes := response.ClientResponse(http.StatusInternalServerError, "could not retrive products by prefix search", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "Successfully retrived all details", productDetails, nil)
	c.JSON(http.StatusOK, successRes)
}

// FilterCategory filters products by category ID.
// @Summary Filter products by category ID
// @Description Filters products based on the provided category ID
// @Tags User Category
// @Accept json
// @Produce json
// @Param category_id query integer true "Category ID"
// @Success 200 {object} response.Response "Success"
// @Failure 400 {object} response.Response "Invalid category ID"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /user/products/filter [get]
func (i *InventoryHandler) FilterCategory(c *gin.Context) {

	CategoryId := c.Query("category_id")
	CategoryIdInt, err := strconv.Atoi(CategoryId)

	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "products Cannot be displayed", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	product_list, err := i.InventoryUseCase.FilterByCategory(CategoryIdInt)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "products cannot be displayed", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	sucessRes := response.ClientResponse(http.StatusOK, "Products List", product_list, nil)
	c.JSON(http.StatusOK, sucessRes)

}

// ProductRating handles the rating of a product by a user.
// @Summary Rate a product
// @Description Allows a user to rate a specific product
// @Tags User Product
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param product_id query integer true "Product ID"
// @Param rating query integer true "Rating (1-5)"
// @Success 200 {object} response.Response "Success"
// @Failure 400 {object} response.Response "Invalid parameters"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /user/products/rating [post]
func (i *InventoryHandler) ProductRating(c *gin.Context) {
	idString, _ := c.Get("id")
	id, _ := idString.(int)

	productID, err := strconv.Atoi(c.Query("product_id"))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "check parameters properly", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	rating, err := strconv.Atoi(c.Query("rating"))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "check parameters properly", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	err = i.InventoryUseCase.ProductRating(id, productID, float64(rating))
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "rating cannot be displayed", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Suceessfully rated the product", nil, nil)
	c.JSON(http.StatusOK, successRes)
}
