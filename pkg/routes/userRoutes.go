package routes

import (
	"HeadZone/pkg/api/handler"
	"HeadZone/pkg/api/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(engine *gin.RouterGroup, userHandler *handler.UserHandler, otpHandler *handler.OtpHandler, inventoryHandler *handler.InventoryHandler, cartHandler *handler.CartHandler, orderHandler *handler.OrderHandler, paymentHandler *handler.PaymentHandler, walletHandler *handler.WalletHandler, couponHandler *handler.CouponHandler) {

	engine.POST("/signup", userHandler.UserSignUp)
	engine.POST("/login", userHandler.LoginHandler)

	engine.POST("/otplogin", otpHandler.SendOTP)
	engine.POST("/verifyotp", otpHandler.VerifyOTP)

	engine.GET("/products", inventoryHandler.ListProducts)

	engine.GET("/payment", paymentHandler.MakePaymentRazorpay) // Update this route
	engine.GET("/verifypayment", paymentHandler.VerifyPayment) // Update this route

	engine.Use(middleware.UserAuthMiddleware)
	{

		profile := engine.Group("/profile")
		{
			profile.GET("", userHandler.GetUserDetails)
			profile.GET("/address", userHandler.GetAddresses)
			profile.POST("", userHandler.AddAddress)
			profile.PUT("", userHandler.EditDetails)
			profile.PATCH("", userHandler.ChangePassword)

			orders := profile.Group("/orders")
			{
				orders.GET("", orderHandler.GetOrders)
				orders.GET("/all", orderHandler.GetAllOrders)
				orders.DELETE("", orderHandler.CancelOrder)
				orders.PUT("/return", orderHandler.ReturnOrder)
			}

		}

		products := engine.Group("/products")
		{
			products.POST("/search", inventoryHandler.SearchProducts)
			products.POST("/filter", inventoryHandler.FilterCategory)
		}

		cart := engine.Group("/cart")
		{
			cart.POST("", cartHandler.AddToCart)
			cart.GET("", userHandler.GetCart)
			cart.DELETE("", userHandler.RemoveFromCart)
			cart.PUT("", userHandler.UpdateQuantity)
		}

		checkout := engine.Group("/check-out")
		{
			checkout.GET("", cartHandler.CheckOut)
			checkout.POST("", orderHandler.OrderItemsFromCart)
			checkout.GET("/print", orderHandler.PrintInvoice)
		}

		wallet := engine.Group("/wallet")
		{
			wallet.GET("", walletHandler.ViewWallet)

		}
		coupon := engine.Group("/coupon")
		{
			coupon.GET("", couponHandler.GetAllCoupons)
		}

	}
}
