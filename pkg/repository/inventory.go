package repository

import (
	"HeadZone/pkg/domain"
	"HeadZone/pkg/repository/interfaces"
	"HeadZone/pkg/utils/models"
	"errors"
	"fmt"
	"strconv"

	"gorm.io/gorm"
)

type inventoryRepository struct {
	DB *gorm.DB
}

func NewInventoryRepository(DB *gorm.DB) interfaces.InventoryRepository {
	return &inventoryRepository{
		DB: DB,
	}
}

func (i *inventoryRepository) AddInventory(inventory models.AddInventories) (models.InventoryResponse, error) {
	// Check if the product already exists based on product name and category ID
	var count int64
	i.DB.Model(&models.Inventory{}).Where("product_name = ? AND category_id = ?", inventory.ProductName, inventory.CategoryID).Count(&count)
	if count > 0 {
		// Product already exists, return an error or handle as needed
		return models.InventoryResponse{}, errors.New("product already exists in the database")
	}

	// Check for negative stock or price values
	if inventory.Stock < 0 || inventory.Price < 0 {
		return models.InventoryResponse{}, errors.New("stock and price cannot be negative")
	}

	// If the product doesn't exist and values are valid, proceed with inserting it into the database
	query := `
        INSERT INTO inventories (category_id, product_name, color, stock, price)
        VALUES (?, ?, ?, ?, ?);
    `
	err := i.DB.Exec(query, inventory.CategoryID, inventory.ProductName, inventory.Color, inventory.Stock, inventory.Price).Error
	if err != nil {
		return models.InventoryResponse{}, err
	}

	var inventoryResponse models.InventoryResponse
	// Populate inventoryResponse as needed

	return inventoryResponse, nil
}

func (prod *inventoryRepository) ListProducts(pageList, offset int) ([]models.InventoryUserResponse, error) {
	var product_list []models.InventoryUserResponse

	query := `
		SELECT i.id, i.category_id, c.category, i.product_name, i.color, i.price 
		FROM inventories i 
		INNER JOIN categories c ON i.category_id = c.id 
		LIMIT $1 OFFSET $2
	`
	fmt.Println(pageList, offset)
	err := prod.DB.Raw(query, pageList, offset).Scan(&product_list).Error

	if err != nil {
		return []models.InventoryUserResponse{}, errors.New("error checking user details")
	}
	fmt.Println("product list", product_list)
	return product_list, nil
}

func (db *inventoryRepository) EditInventory(inventory domain.Inventory, id int) (domain.Inventory, error) {

	var modInventory domain.Inventory

	query := "UPDATE inventories SET category_id = ?, product_name = ?, color = ?, stock = ?, price = ? WHERE id = ?"

	if err := db.DB.Exec(query, inventory.CategoryID, inventory.ProductName, inventory.Color, inventory.Stock, inventory.Price, id).Error; err != nil {
		return domain.Inventory{}, err
	}

	if err := db.DB.First(&modInventory, id).Error; err != nil {
		return domain.Inventory{}, err
	}
	return modInventory, nil
}

func (i *inventoryRepository) DeleteInventory(inventoryID string) error {

	id, err := strconv.Atoi(inventoryID)
	if err != nil {
		return errors.New("converting into integet is not happened")
	}

	result := i.DB.Exec("DELETE FROM inventories WHERE id = ?", id)

	if result.RowsAffected < 1 {
		return errors.New("no records with that ID exist")
	}

	return nil
}

func (i *inventoryRepository) CheckInventory(pid int) (bool, error) {
	var k int
	err := i.DB.Raw("SELECT COUNT(*) FROM inventories WHERE id=?", pid).Scan(&k).Error
	if err != nil {
		return false, err
	}

	if k == 0 {
		return false, err
	}

	return true, err
}

func (i *inventoryRepository) UpdateInventory(pid int, stock int) (models.InventoryResponse, error) {

	// Check the database connection
	if i.DB == nil {
		return models.InventoryResponse{}, errors.New("database connection is nil")
	}

	// Update the
	if err := i.DB.Exec("UPDATE inventories SET stock = stock + $1 WHERE id= $2", stock, pid).Error; err != nil {
		return models.InventoryResponse{}, err
	}

	// Retrieve the update
	var newdetails models.InventoryResponse
	var newstock int
	if err := i.DB.Raw("SELECT stock FROM inventories WHERE id=?", pid).Scan(&newstock).Error; err != nil {
		return models.InventoryResponse{}, err
	}
	newdetails.ProductID = pid
	newdetails.Stock = newstock

	return newdetails, nil
}

func (i *inventoryRepository) ShowIndividualProducts(id string) (models.InventoryUserResponse, error) {
	pid, error := strconv.Atoi(id)
	if error != nil {
		return models.InventoryUserResponse{}, errors.New("convertion not happened")
	}
	var product models.InventoryUserResponse
	err := i.DB.Raw(`
	SELECT
		*
		FROM
			inventories
		
		WHERE
			inventories.id = ?
			`, pid).Scan(&product).Error

	if err != nil {
		return models.InventoryUserResponse{}, errors.New("error retrieved record")
	}
	return product, nil

}

func (i *inventoryRepository) CheckStock(pid int) (int, error) {
	var k int
	if err := i.DB.Raw("SELECT stock FROM inventories WHERE id=$1", pid).Scan(&k).Error; err != nil {
		return 0, err
	}
	return k, nil

}

func (c *inventoryRepository) FetchProductDetails(productId uint) (models.Inventory, error) {
	var product models.Inventory
	err := c.DB.Raw(`SELECT price,stock FROM inventories WHERE id=?`, productId).Scan(&product).Error
	return product, err
}
