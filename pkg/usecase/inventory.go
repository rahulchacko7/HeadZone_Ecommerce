package usecase

import (
	"HeadZone/pkg/domain"
	"HeadZone/pkg/helper/interfaces"
	repo "HeadZone/pkg/repository/interfaces"
	usecase "HeadZone/pkg/usecase/interfaces"
	"HeadZone/pkg/utils/models"
	"errors"
	"strings"
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

func (i *inventoryUseCase) AddInventory(inventory models.AddInventories) (models.Inventory, error) {

	if inventory.Stock < 0 || inventory.Price < 0 || inventory.CategoryID < 0 {
		return models.Inventory{}, errors.New("negative values not allowed for stock, price, or category ID")
	}

	inventoryResponse, err := i.repository.AddInventory(inventory)
	if err != nil {
		return models.Inventory{}, err
	}

	return inventoryResponse, nil
}

func (i *inventoryUseCase) ListProducts(pageNo, pageList int) ([]models.InventoryUserResponse, error) {

	if pageList <= 0 || pageNo <= 0 {
		return []models.InventoryUserResponse{}, errors.New("there is no inventory as you mentioned")
	}

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

	if pid <= 0 || stock <= 0 {
		return models.InventoryResponse{}, errors.New("must enter a positive value")
	}

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

func (i *inventoryUseCase) SearchProductsOnPrefix(prefix string) ([]models.InventoryUserResponse, error) {

	inventoryList, err := i.repository.GetInventory(prefix)

	if err != nil {
		return nil, err
	}

	var filteredProducts []models.InventoryUserResponse

	for _, product := range inventoryList {
		if strings.HasPrefix(strings.ToLower(product.ProductName), strings.ToLower(prefix)) {
			filteredProducts = append(filteredProducts, product)
		}
	}

	if len(filteredProducts) == 0 {
		return nil, errors.New("no items matching your keyword")
	}

	return filteredProducts, nil
}

func (i *inventoryUseCase) FilterByCategory(CategoryIdInt int) ([]models.InventoryUserResponse, error) {

	if CategoryIdInt <= 0 {
		return []models.InventoryUserResponse{}, errors.New("id must be a positive value")
	}

	product_list, err := i.repository.FilterByCategory(CategoryIdInt)

	if err != nil {
		return []models.InventoryUserResponse{}, err
	}

	return product_list, nil
}
