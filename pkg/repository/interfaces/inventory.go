package interfaces

import (
	"HeadZone/pkg/domain"
	"HeadZone/pkg/utils/models"
)

type InventoryRepository interface {
	AddInventory(inventory models.AddInventories) (models.InventoryResponse, error)
	ListProducts(int, int) ([]models.Inventory, error)
	EditInventory(domain.Inventory, int) (domain.Inventory, error)
	DeleteInventory(id string) error
}
