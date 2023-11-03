package routes

import (
	"HeadZone/pkg/api/handler"
	"HeadZone/pkg/api/middleware"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(engine *gin.RouterGroup, adminHandler *handler.AdminHandler, categoryHandler *handler.CategoryHandler) {

	engine.POST("/adminlogin", adminHandler.LoginHandler)

	engine.Use(middleware.AdminAuthMiddleware)
	{
		usermanagement := engine.Group("/users")
		{
			usermanagement.POST("/block", adminHandler.BlockUser)
			usermanagement.POST("/unblock", adminHandler.UnBlockUser)
			usermanagement.GET("/", adminHandler.GetUsers)
		}

		categorymanagement := engine.Group("/category")
		{
			categorymanagement.POST("/add", categoryHandler.AddCategory)
			categorymanagement.GET("/getcategory", categoryHandler.GetCategory)
			categorymanagement.PUT("/updatecategory", categoryHandler.UpdateCategory)
			// categorymanagement.DELETE("/deletecategory", categoryHandler.DeleteCategory)
		}
	}
}
