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

func (i *orderRepository) GetAllOrders(userID, page, pageSize int) ([]models.OrderDetails, error) {
	if page == 0 {
		page = 1
	}
	offset := (page - 1) * pageSize
	var order []models.OrderDetails

	err := i.DB.Raw("SELECT id as order_id, address_id, payment_method_id, final_price as price, order_status, payment_status FROM orders WHERE user_id = ? OFFSET ? LIMIT ?", userID, offset, pageSize).Scan(&order).Error
	if err != nil {
		return nil, err
	}
	fmt.Println("Retrieved orders:", order)
	return order, nil
}
