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
func (i *InventoryHandler) AddInventory(c *gin.Context) {
	var inventory models.AddInventories

	if err := c.ShouldBindJSON(&inventory); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "form file error", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	fmt.Println("kkkkkkkkkk", inventory)
	InventoryResponse, err := i.InventoryUseCase.AddInventory(inventory)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "Could not add the Inventory", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully added inventory", InventoryResponse, nil)
	c.JSON(http.StatusOK, successRes)

}

// ListProducts handles the retrieval of products with pagination.
// @Summary List products with pagination
// @Description Retrieves a list of products with pagination support
// @Tags products
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
