package routes

import (
	"HeadZone/pkg/api/handler"
	"HeadZone/pkg/api/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(engine *gin.RouterGroup, userHandler *handler.UserHandler, otpHandler *handler.OtpHandler, inventoryHandler *handler.InventoryHandler, cartHandler *handler.CartHandler) {

	engine.POST("/signup", userHandler.UserSignUp)
	engine.POST("/login", userHandler.LoginHandler)

	engine.POST("/otplogin", otpHandler.SendOTP)
	engine.POST("/verifyotp", otpHandler.VerifyOTP)

	engine.GET("/viewproducts", inventoryHandler.ListProducts)

	engine.Use(middleware.UserAuthMiddleware)
	{

		profile := engine.Group("/profile")
		{
			profile.GET("/userdetails", userHandler.GetUserDetails)
			profile.GET("/view", userHandler.GetAddresses)
			profile.POST("/address/add", userHandler.AddAddress)

			edit := profile.Group("/edit")
			{
				edit.PUT("/name", userHandler.EditName)
				edit.PUT("/email", userHandler.EditEmail)
				edit.PUT("/phone", userHandler.EditPhone)
			}

			security := profile.Group("/security")
			{
				security.PUT("/change-password", userHandler.ChangePassword)
			}

		}

		home := engine.Group("/home")
		{
			home.GET("/products", inventoryHandler.ListProducts)
			home.GET("/product/details", inventoryHandler.ShowIndividualProducts)
			home.POST("/add-to-cart", cartHandler.AddToCart)
			//home.POST("/whishlist/add", whislisthandler.AddWishlist)
		}

	}
}
