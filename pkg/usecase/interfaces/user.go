package interfaces

import "HeadZone/pkg/utils/models"

type UserUseCase interface {
	UserSignUp(models.UserDetails) (models.TokenUsers, error)
	LoginHandler(models.UserLogin) (models.TokenUsers, error)
}
