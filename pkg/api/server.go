package http

import (
	"HeadZone/pkg/api/handler"
	"HeadZone/pkg/routes"
	"log"

	"github.com/gin-gonic/gin"
)

type ServerHTTP struct {
	engine *gin.Engine
}

func NewServerHTTP(userHandler *handler.UserHandler, adminHandler *handler.AdminHandler, otpHandler *handler.OtpHandler, categoryHandler *handler.CategoryHandler, inventoryHandler *handler.InventoryHandler, cartHandler *handler.CartHandler, orderHandler *handler.OrderHandler, paymentHandler *handler.PaymentHandler, walletHandler *handler.WalletHandler) *ServerHTTP {
	engine := gin.New()

	engine.Use(gin.Logger())
	engine.LoadHTMLGlob("pkg/templates/index.html")
	engine.GET("/validate-token", adminHandler.ValidateRefreshTokenAndCreateNewAccess)

	routes.UserRoutes(engine.Group("/user"), userHandler, otpHandler, inventoryHandler, cartHandler, orderHandler, paymentHandler, walletHandler)
	routes.AdminRoutes(engine.Group("/admin"), adminHandler, categoryHandler, inventoryHandler, orderHandler)

	return &ServerHTTP{engine: engine}
}

func (sh *ServerHTTP) Start() {
	err := sh.engine.Run(":3000")
	if err != nil {
		log.Fatal("gin engine couldn't start")
	}
}
