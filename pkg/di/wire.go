//go:build wireinject
// +build wireinject

package di

import (
	http "HeadZone/pkg/api"
	"HeadZone/pkg/api/handler"
	config "HeadZone/pkg/config"
	db "HeadZone/pkg/db"
	repository "HeadZone/pkg/repository"
	usecase "HeadZone/pkg/usecase"

	"github.com/google/wire"
)

func InitializeAPI(cfg config.Config) (*http.ServerHTTP, error) {
	wire.Build(db.ConnectDatabase, repository.NewUserRepository, usecase.NewUserUseCase, handler.NewUserHandler, http.NewServerHTTP)

	return &http.ServerHTTP{}, nil
}
