package routes

import (
	"HeadZone/pkg/api/handler"
	"HeadZone/pkg/api/middleware"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(engine *gin.RouterGroup, adminHandler *handler.AdminHandler) {

	engine.POST("/adminlogin", adminHandler.LoginHandler)

	engine.Use(middleware.AdminAuthMiddleware)
	{
		usermanagement := engine.Group("/users")
		{
			usermanagement.POST("/block", adminHandler.BlockUser)
			usermanagement.POST("/unblock", adminHandler.UnBlockUser)
			usermanagement.GET("/", adminHandler.GetUsers)
		}
	}
}
