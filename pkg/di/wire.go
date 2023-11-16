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
		repository.NewAdminRepository,
		repository.NewCategoryRepository,
		repository.NewOtpRepository,
		repository.NewInventoryRepository,
		repository.NewCartRepository,
		repository.NewUserRepository,
		helper.NewHelper,

		usecase.NewUserUseCase,
		usecase.NewAdminUseCase,
		usecase.NewCategoryUseCase,
		usecase.NewInventoryUseCase,
		usecase.NewOtpUseCase,
		usecase.NewCartUseCase,

		handler.NewUserHandler,
		handler.NewAdminHandler,
		handler.NewInventoryHandler,
		handler.NewOtpHandler,
		handler.NewCategoryHandler,
		handler.NewCartHandler,

		http.NewServerHTTP,
	)

	return &http.ServerHTTP{}, nil
}
