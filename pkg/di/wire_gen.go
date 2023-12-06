// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package di

import (
	"HeadZone/pkg/api"
	"HeadZone/pkg/api/handler"
	"HeadZone/pkg/config"
	"HeadZone/pkg/db"
	"HeadZone/pkg/helper"
	"HeadZone/pkg/repository"
	"HeadZone/pkg/usecase"
)

// Injectors from wire.go:

func InitializeAPI(cfg config.Config) (*http.ServerHTTP, error) {
	gormDB, err := db.ConnectDatabase(cfg)
	if err != nil {
		return nil, err
	}
	userRepository := repository.NewUserRepository(gormDB)
	otpRepository := repository.NewOtpRepository(gormDB)
	inventoryRepository := repository.NewInventoryRepository(gormDB)
	interfacesHelper := helper.NewHelper(cfg)
	userUseCase := usecase.NewUserUseCase(userRepository, cfg, otpRepository, inventoryRepository, interfacesHelper)
	userHandler := handler.NewUserHandler(userUseCase)
	adminRepository := repository.NewAdminRepository(gormDB)
	adminUseCase := usecase.NewAdminUseCase(adminRepository, interfacesHelper)
	adminHandler := handler.NewAdminHandler(adminUseCase)
	otpUseCase := usecase.NewOtpUseCase(cfg, otpRepository, interfacesHelper)
	otpHandler := handler.NewOtpHandler(otpUseCase)
	categoryRepository := repository.NewCategoryRepository(gormDB)
	categoryUseCase := usecase.NewCategoryUseCase(categoryRepository, inventoryRepository)
	categoryHandler := handler.NewCategoryHandler(categoryUseCase)
	inventoryUseCase := usecase.NewInventoryUseCase(inventoryRepository, interfacesHelper)
	inventoryHandler := handler.NewInventoryHandler(inventoryUseCase)
	cartRepository := repository.NewCartRepository(gormDB)
	cartUseCase := usecase.NewCartUseCase(cartRepository, inventoryRepository, userUseCase, adminRepository)
	cartHandler := handler.NewCartHandler(cartUseCase)
	orderRepository := repository.NewOrderRepository(gormDB)
	walletRepository := repository.NewWalletRepository(gormDB)
	orderUseCase := usecase.NewOrderUseCase(orderRepository, userUseCase, walletRepository, cartRepository)
	orderHandler := handler.NewOrderHandler(orderUseCase)
	paymentRepository := repository.NewPaymentRepository(gormDB)
	paymentUseCase := usecase.NewPaymentUseCase(orderRepository, paymentRepository)
	paymentHandler := handler.NewPaymentHandler(paymentUseCase)
	walletUsecase := usecase.NewWalletUseCase(walletRepository)
	walletHandler := handler.NewWalletHandler(walletUsecase)
	couponRepository := repository.NewCouponRepository(gormDB)
	couponUseCase := usecase.NewCouponUseCase(couponRepository)
	couponHandler := handler.NewCouponHandler(couponUseCase)
	serverHTTP := http.NewServerHTTP(userHandler, adminHandler, otpHandler, categoryHandler, inventoryHandler, cartHandler, orderHandler, paymentHandler, walletHandler, couponHandler)
	return serverHTTP, nil
}
