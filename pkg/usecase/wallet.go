package usecase

import (
	"HeadZone/pkg/repository/interfaces"
	services "HeadZone/pkg/usecase/interfaces"

	"HeadZone/pkg/utils/models"
)

type walletUseCase struct {
	walletRepository interfaces.WalletRepository
}

func NewWalletUseCase(repository interfaces.WalletRepository) services.WalletUsecase {
	return &walletUseCase{
		walletRepository: repository,
	}
}

func (wt *walletUseCase) GetWallet(userID int) (models.WalletAmount, error) {
	return wt.walletRepository.GetWallet(userID)
}
