package interfaces

import (
	"HeadZone/pkg/domain"
	"HeadZone/pkg/utils/models"
)

type InventoryUseCase interface {
	AddInventory(inventory models.AddInventories) (models.InventoryResponse, error)
	ListProducts(int, int) ([]models.InventoryUserResponse, error)
	EditInventory(domain.Inventory, int) (domain.Inventory, error)
	DeleteInventory(id string) error
	UpdateInventory(productID int, stock int) (models.InventoryResponse, error)
	ShowIndividualProducts(id string) (models.InventoryUserResponse, error)
}
