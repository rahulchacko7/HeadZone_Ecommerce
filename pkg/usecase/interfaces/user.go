package interfaces

import "HeadZone/pkg/utils/models"

type UserUseCase interface {
	UserSignUp(user models.UserDetails) (models.TokenUsers, error)
	LoginHandler(user models.UserLogin) (models.TokenUsers, error)
}
