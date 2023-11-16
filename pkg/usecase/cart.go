package usecase

import (
	interfaces "HeadZone/pkg/repository/interfaces"
	"errors"
)

type cartUseCase struct {
	repo                interfaces.CartRepository
	inventoryRepository interfaces.InventoryRepository
	// userUseCase         services.UserUseCase
}

func NewCartUseCase(repo interfaces.CartRepository, inventoryRepo interfaces.InventoryRepository) *cartUseCase {
	return &cartUseCase{
		repo:                repo,
		inventoryRepository: inventoryRepo,
		// userUseCase:         userUseCase,
	}
}

func (i *cartUseCase) AddToCart(userID, inventoryID int) error {

	//check if item already added if already present send error as already added

	//check if the desired product has quantity available
	stock, err := i.inventoryRepository.CheckStock(inventoryID)
	if err != nil {
		return err
	}
	//if available then call userRepository
	if stock <= 0 {
		return errors.New("out of stock")
	}

	//find user cart id
	cart_id, err := i.repo.GetCartId(userID)
	if err != nil {
		return errors.New("some error in geting user cart")
	}
	//if user has no existing cart create new cart
	if cart_id == 0 {
		cart_id, err = i.repo.CreateNewCart(userID)
		if err != nil {
			return errors.New("cannot create cart fro user")
		}
	}

	exists, err := i.repo.CheckIfItemIsAlreadyAdded(cart_id, inventoryID)
	if err != nil {
		return err
	}

	if exists {
		return errors.New("item already exists in cart")
	}

	//add product to line items
	if err := i.repo.AddLineItems(cart_id, inventoryID); err != nil {
		return errors.New("error in adding products")
	}

	return nil
}
