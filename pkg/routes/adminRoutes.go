package routes

import (
	"HeadZone/pkg/api/handler"
	"HeadZone/pkg/api/middleware"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(engine *gin.RouterGroup, adminHandler *handler.AdminHandler, categoryHandler *handler.CategoryHandler, inventoryHandler *handler.InventoryHandler) {

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
			categorymanagement.POST("", categoryHandler.AddCategory)
			categorymanagement.GET("", categoryHandler.GetCategory)
			categorymanagement.PUT("", categoryHandler.UpdateCategory)
			categorymanagement.DELETE("", categoryHandler.DeleteCategory)
		}

		inventorymanagement := engine.Group("/inventory")
		{
			inventorymanagement.POST("", inventoryHandler.AddInventory)
			inventorymanagement.GET("", inventoryHandler.ListProducts)
			//inventorymanagement.GET("", inventoryHandler.ViewProductsByID)
			inventorymanagement.PUT("", inventoryHandler.EditInventory)
			inventorymanagement.DELETE("", inventoryHandler.DeleteInventory)
		}
	}
}
