package interfaces

import (
	"HeadZone/pkg/domain"
	"HeadZone/pkg/utils/models"
)

type AdminUseCase interface {
	LoginHandler(adminDetails models.AdminLogin) (domain.TokenAdmin, error)
	BlockUser(id string) error
	UnBlockUser(id string) error
	GetUsers(page int, count int) ([]models.UserDetailsAtAdmin, error)
}
