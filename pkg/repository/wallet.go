package repository

import (
	"HeadZone/pkg/repository/interfaces"
	"HeadZone/pkg/utils/models"
	"fmt"

	"gorm.io/gorm"
)

type walletRepository struct {
	DB *gorm.DB
}

func NewWalletRepository(DB *gorm.DB) interfaces.WalletRepository {
	return &walletRepository{
		DB: DB,
	}
}
func (wt *walletRepository) GetWallet(userID int) (models.WalletAmount, error) {
	var walletAmount models.WalletAmount
	err := wt.DB.Raw("SELECT amount FROM wallets WHERE user_id = ?", userID).Scan(&walletAmount).Error
	if err != nil {
		return models.WalletAmount{}, err
	}
	return walletAmount, nil
}

// func (wt *walletRepository) GetsWallet(orderId int) (models.WalletAmount, error) {
// 	var walletAmount models.WalletAmount
// 	err := wt.DB.Raw("SELECT amount FROM wallets WHERE user_id = ?", orderId).Scan(&walletAmount).Error
// 	if err != nil {
// 		return models.WalletAmount{}, err
// 	}
// 	return walletAmount, nil
// }

func (wt *walletRepository) AddToWallet(Price, UserId int) (models.WalletAmount, error) {
	fmt.Println("amount inside repo", Price)
	fmt.Println("UserId inside repo", UserId)

	var walletAmount models.WalletAmount
	err := wt.DB.Exec("INSERT INTO wallets (user_id,amount) VALUES(?,?)", UserId, 0).Error
	if err != nil {
		return models.WalletAmount{}, err
	}
	err = wt.DB.Exec("UPDATE wallets SET amount = amount + ? where user_id = ?", Price, UserId).Error
	if err != nil {
		return models.WalletAmount{}, err
	}
	return walletAmount, nil
}
