package repository

import (
	"HeadZone/pkg/domain"
	"HeadZone/pkg/repository/interfaces"
	"HeadZone/pkg/utils/models"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type orderRepository struct {
	DB *gorm.DB
}

func NewOrderRepository(db *gorm.DB) interfaces.OrderRepository {
	return &orderRepository{
		DB: db,
	}
}

func (i *orderRepository) OrderItems(userid, addressid, paymentid int, total float64) (int, error) {

	var id int
	query := `
    INSERT INTO orders (created_at,user_id,address_id, payment_method_id, final_price)
    VALUES (Now(),?, ?, ?, ?)
    RETURNING id
    `
	i.DB.Raw(query, userid, addressid, paymentid, total).Scan(&id)
	fmt.Println("id...........", id)
	return id, nil

}

func (i *orderRepository) AddOrderProducts(order_id int, cart []models.GetCart) error {
	query := `
    INSERT INTO order_items (order_id,inventory_id,quantity,total_price)
    VALUES (?, ?, ?, ?)
    `

	for _, v := range cart {
		var inv int
		if err := i.DB.Raw("select id from inventories where product_name=$1", v.ProductName).Scan(&inv).Error; err != nil {
			return err
		}

		if err := i.DB.Exec(query, order_id, inv, v.Quantity, v.Total).Error; err != nil {
			return err
		}
	}

	return nil

}

func (i *orderRepository) GetOrders(orderID int) (domain.OrderResponse, error) {
	if orderID <= 0 {
		return domain.OrderResponse{}, errors.New("order ID should be a positive number")
	}

	var order domain.OrderResponse

	query := `SELECT * FROM orders WHERE id = $1`

	if err := i.DB.Raw(query, orderID).First(&order).Error; err != nil {
		return domain.OrderResponse{}, err
	}

	return order, nil
}

func (o *orderRepository) CheckOrderStatusByID(id int) (string, error) {

	var status string
	err := o.DB.Raw("select order_status from orders where id = ?", id).Scan(&status).Error
	if err != nil {
		return "", err
	}

	return status, nil
}

func (i *orderRepository) CancelOrder(id int) error {

	if err := i.DB.Exec("update orders set order_status='CANCELED' where id=$1", id).Error; err != nil {
		return err
	}

	return nil

}

// func (i *orderRepository) GetOrderDetails(userID, page, count int) ([]models.AllOrderResponse, error) {
// 	if userID <= 0 || page <= 0 || count <= 0 {
// 		return nil, errors.New("please provide positive numbers only")
// 	}

// 	if page == 0 {
// 		page = 1
// 	}
// 	offset := (page - 1) * count

// 	var orders []models.AllOrderResponse

// 	// Fetch orders for the given user within the provided pagination limits
// 	i.DB.Raw("SELECT id as order_id, final_price, shipment_status, payment_status FROM orders WHERE user_id = ? LIMIT ? OFFSET ?", userID, count, offset).Scan(&orders)

// 	// Iterate over the fetched orders and retrieve their associated order items
// 	for idx := range orders {
// 		var orderItems []models.AllOrderResponse
// 		i.DB.Raw(`SELECT
// 			order_items.order_id,
// 			order_items.inventory_id,
// 			order_items.quantity,
// 			order_items.total_price
// 		FROM
// 			order_items
// 		WHERE
// 			order_items.order_id = ?`, orders[idx].OrderDetails.ID).Scan(&orderItems)

// 		orders[idx].AllOrderResponse = orderItems
// 	}

// 	return orders, nil
// }
