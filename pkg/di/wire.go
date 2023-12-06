//go:build wireinject
// +build wireinject

package di

import (
	http "HeadZone/pkg/api"
	"HeadZone/pkg/api/handler"
	config "HeadZone/pkg/config"
	db "HeadZone/pkg/db"
	"HeadZone/pkg/helper"
	repository "HeadZone/pkg/repository"
	usecase "HeadZone/pkg/usecase"

	"github.com/google/wire"
)

func InitializeAPI(cfg config.Config) (*http.ServerHTTP, error) {
	wire.Build(
		db.ConnectDatabase,
		repository.NewUserRepository,
		usecase.NewUserUseCase,
		handler.NewUserHandler,
		helper.NewHelper,

		handler.NewAdminHandler,
		handler.NewInventoryHandler,
		handler.NewOtpHandler,
		handler.NewCategoryHandler,
		handler.NewCartHandler,
		handler.NewOrderHandler,
		handler.NewPaymentHandler,
		handler.NewWalletHandler,
		handler.NewCouponHandler,

		usecase.NewAdminUseCase,
		usecase.NewCategoryUseCase,
		usecase.NewInventoryUseCase,
		usecase.NewOtpUseCase,
		usecase.NewCartUseCase,
		usecase.NewOrderUseCase,
		usecase.NewPaymentUseCase,
		usecase.NewWalletUseCase,
		usecase.NewCouponUseCase,

		repository.NewAdminRepository,
		repository.NewCategoryRepository,
		repository.NewOtpRepository,
		repository.NewInventoryRepository,
		repository.NewCartRepository,
		repository.NewOrderRepository,
		repository.NewPaymentRepository,
		repository.NewWalletRepository,
		repository.NewCouponRepository,

		http.NewServerHTTP,
	)

	return &http.ServerHTTP{}, nil
}
