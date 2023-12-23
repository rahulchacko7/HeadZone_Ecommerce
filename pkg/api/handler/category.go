package handler

import (
	"HeadZone/pkg/domain"
	"HeadZone/pkg/usecase/interfaces"
	"HeadZone/pkg/utils/models"
	"HeadZone/pkg/utils/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	CategoryUseCase interfaces.CategoryUseCase
}

func NewCategoryHandler(usecase interfaces.CategoryUseCase) *CategoryHandler {
	return &CategoryHandler{
		CategoryUseCase: usecase,
	}
}

// AddCategory handles the addition of a new category.
// @Summary Add a new category
// @Description Add a new category to the system
// @Tags Admin Category Management
// @Accept json
// @Produce json
// @security BearerTokenAuth
// @Param body body domain.Category true "Category object to be added"
// @Success 200 {object} domain.Category "Successfully added category"
// @Failure 400 {object} response.Response "Invalid request or incorrect format"
// @Router /admin/category [post]
func (cat *CategoryHandler) AddCategory(c *gin.Context) {

	var category domain.Category
	if err := c.BindJSON(&category); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "fields are provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	CategoryResponse, err := cat.CategoryUseCase.AddCategory(category)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "Could not add the category", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Sucessfully added Category", CategoryResponse, nil)
	c.JSON(http.StatusOK, successRes)
}

// @Tags Admin Category Management
func (Cat *CategoryHandler) GetCategory(c *gin.Context) {

	categories, err := Cat.CategoryUseCase.GetCategories()
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully got all categories", categories, nil)
	c.JSON(http.StatusOK, successRes)

}

// UpdateCategory handles the updating of a category name.
// @Summary Update a category name
// @Description Update an existing category name in the system
// @Tags Admin Category Management
// @Accept json
// @Produce json
// @security BearerTokenAuth
// @Param request body models.SetNewName true "Current and New names in JSON format"
// @Success 200 {object} domain.Category "Successfully updated category name"
// @Failure 400 {object} response.Response "Invalid request or incorrect format"
// @Failure 400 {object} response.Response "Could not update the category"
// @Router /admin/category [put]
func (Cat *CategoryHandler) UpdateCategory(c *gin.Context) {
	var p models.SetNewName

	if err := c.BindJSON(&p); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "Fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	a, err := Cat.CategoryUseCase.UpdateCategory(p.Current, p.New)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not update the category", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Sucessfully updated...", a, nil)
	c.JSON(http.StatusOK, successRes)
}

// @Tags Admin Category Management
func (Cat *CategoryHandler) DeleteCategory(c *gin.Context) {

	categoryID := c.Query("id")
	err := Cat.CategoryUseCase.DeleteCategory(categoryID)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "Fields are not provided in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	SuccessRes := response.ClientResponse(http.StatusOK, "Sucessfully Deleted...", nil, nil)
	c.JSON(http.StatusOK, SuccessRes)
}
