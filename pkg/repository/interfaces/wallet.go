package interfaces

import "HeadZone/pkg/utils/models"

type WalletRepository interface {
	GetWallet(userID int) (models.WalletAmount, error)
	GetsWallet(orderId int) (models.WalletAmount, error)
	AddToWallet(Price, UserId int) (models.WalletAmount, error)
}
