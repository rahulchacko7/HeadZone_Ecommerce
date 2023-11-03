package di

import (
	http "HeadZone/pkg/api"
	"HeadZone/pkg/api/handler"
	"HeadZone/pkg/config"
	"HeadZone/pkg/db"
	"HeadZone/pkg/repository"
	"HeadZone/pkg/usecase"
)

func InitializeAPI(cfg config.Config) (*http.ServerHTTP, error) {
	gormDB, err := db.ConnectDatabase(cfg)
	if err != nil {
		return nil, err
	}

	adminRepository := repository.NewAdminRepository(gormDB)
	adminUseCase := usecase.NewAdminUseCase(adminRepository)
	adminHandler := handler.NewAdminHandler(adminUseCase)

	otpRepository := repository.NewOtpRepository(gormDB)
	otpUseCase := usecase.NewOtpUseCase(cfg, otpRepository)
	otpHandler := handler.NewOtpHandler(otpUseCase)

	userRepository := repository.NewUserRepository(gormDB)
	userUseCase := usecase.NewUserUseCase(userRepository, cfg, otpRepository)
	userHandler := handler.NewUserHandler(userUseCase)

	categoryRespository := repository.NewCategoryRepository(gormDB)
	categoryUseCase := usecase.NewCategoryUseCase(categoryRespository)
	categoryHandler := handler.NewCategoryHandler(categoryUseCase)

	serverHTTP := http.NewServerHTTP(userHandler, adminHandler, otpHandler, categoryHandler)

	return serverHTTP, nil
}
