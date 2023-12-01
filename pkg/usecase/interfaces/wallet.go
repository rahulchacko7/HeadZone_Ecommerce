package interfaces

import "HeadZone/pkg/utils/models"

type WalletUsecase interface {
	GetWallet(id int) (models.WalletAmount, error)
}
