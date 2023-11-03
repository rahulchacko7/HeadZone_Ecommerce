package http

import (
	"HeadZone/pkg/api/handler"
	"HeadZone/pkg/routes"

	"github.com/gin-gonic/gin"
)

type ServerHTTP struct {
	engine *gin.Engine
}

func NewServerHTTP(userHandler *handler.UserHandler, adminHandler *handler.AdminHandler, otpHandler *handler.OtpHandler, categoryHandler *handler.CategoryHandler) *ServerHTTP {
	engine := gin.New()

	engine.Use(gin.Logger())

	routes.UserRoutes(engine.Group("/user"), userHandler, otpHandler)
	routes.AdminRoutes(engine.Group("/admin"), adminHandler, categoryHandler)

	return &ServerHTTP{engine: engine}
}

func (sh *ServerHTTP) Start() {
	sh.engine.Run(":3000")
}
