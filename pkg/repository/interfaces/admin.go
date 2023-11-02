package interfaces

import (
	"HeadZone/pkg/domain"
	"HeadZone/pkg/utils/models"
)

type AdminRepository interface {
	LoginHandler(adminDetails models.AdminLogin) (domain.Admin, error)
	GetUserByID(id string) (domain.Users, error)
	UpdateBlockUserByID(user domain.Users) error
	GetUsers(page int, count int) ([]models.UserDetailsAtAdmin, error)
}
