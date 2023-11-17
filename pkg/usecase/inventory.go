package usecase

import (
	"HeadZone/pkg/domain"
	"HeadZone/pkg/helper/interfaces"
	repo "HeadZone/pkg/repository/interfaces"
	usecase "HeadZone/pkg/usecase/interfaces"
	"HeadZone/pkg/utils/models"
	"errors"
)

type inventoryUseCase struct {
	repository repo.InventoryRepository
	helper     interfaces.Helper
}

func NewInventoryUseCase(repo repo.InventoryRepository, h interfaces.Helper) usecase.InventoryUseCase {
	return &inventoryUseCase{
		repository: repo,
		helper:     h,
	}
}

func (i *inventoryUseCase) AddInventory(inventory models.AddInventories) (models.InventoryResponse, error) {

	// url, err := i.helper.AddImageToS3(image)
	// if err != nil {
	// 	return models.InventoryResponse{}, err
	// }

	//send the url and save it in database
	InventoryResponse, err := i.repository.AddInventory(inventory)
	if err != nil {
		return models.InventoryResponse{}, err
	}

	return InventoryResponse, nil

}

func (i *inventoryUseCase) ListProducts(pageNo, pageList int) ([]models.InventoryUserResponse, error) {

	offset := (pageNo - 1) * pageList
	productList, err := i.repository.ListProducts(pageList, offset)
	if err != nil {
		return []models.InventoryUserResponse{}, err
	}
	return productList, nil
}

func (usecase *inventoryUseCase) EditInventory(inventory domain.Inventory, id int) (domain.Inventory, error) {
	modInventory, err := usecase.repository.EditInventory(inventory, id)
	if err != nil {
		return domain.Inventory{}, err
	}
	return modInventory, nil
}

func (usecase *inventoryUseCase) DeleteInventory(inventoryID string) error {

	err := usecase.repository.DeleteInventory(inventoryID)
	if err != nil {
		return err
	}
	return nil
}

func (i inventoryUseCase) UpdateInventory(pid int, stock int) (models.InventoryResponse, error) {

	result, err := i.repository.CheckInventory(pid)
	if err != nil {
		return models.InventoryResponse{}, err
	}

	if !result {
		return models.InventoryResponse{}, errors.New("there is no inventory as you mentioned")
	}

	newcat, err := i.repository.UpdateInventory(pid, stock)
	if err != nil {
		return models.InventoryResponse{}, err
	}

	return newcat, err
}

func (i *inventoryUseCase) ShowIndividualProducts(id string) (models.InventoryUserResponse, error) {

	product, err := i.repository.ShowIndividualProducts(id)
	if err != nil {
		return models.InventoryUserResponse{}, err
	}

	return product, nil
}
