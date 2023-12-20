package main

import (
	_ "HeadZone/cmd/docs"
	"HeadZone/pkg/config"
	"HeadZone/pkg/di"

	"log"

	_ "github.com/swaggo/files"
	_ "github.com/swaggo/gin-swagger"
)

// @title Go + Gin E-Commerce API
// @version 1.0.0
// @description TechDeck is an E-commerce platform to purchase and sell Electronic itmes
// @contact.name API Support
// @securityDefinitions.apikey BearerTokenAuth
// @in header
// @name Authorization
// @host localhost:3000
// @BasePath /
// @query.collection.format multi

func main() {

	config, configErr := config.LoadConfig()
	if configErr != nil {
		log.Fatal("cannot load config: ", configErr)
	}
	server, diErr := di.InitializeAPI(config)
	if diErr != nil {
		log.Fatal("cannot start server: ", diErr)
	} else {
		server.Start()
	}
}
