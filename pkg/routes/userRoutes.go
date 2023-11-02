package routes

import (
	"HeadZone/pkg/api/handler"

	"github.com/gin-gonic/gin"
)

func UserRoutes(engine *gin.RouterGroup, userHandler *handler.UserHandler, otpHandler *handler.OtpHandler) {

	engine.POST("/signup", userHandler.UserSignUp)
	engine.POST("/login", userHandler.LoginHandler)

	engine.POST("/otplogin", otpHandler.SendOTP)
	engine.POST("/verifyotp", otpHandler.VerifyOTP)
}
