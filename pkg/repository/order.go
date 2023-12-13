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

func (i *orderRepository) ReduceInventoryQuantity(productName string, quantity int) error {
	query := `
        UPDATE inventories
        SET stock = stock - ?
        WHERE product_name = ?
    `
	if err := i.DB.Exec(query, quantity, productName).Error; err != nil {
		return err
	}
	return nil
}

func (i *orderRepository) GetOrders(orderID int) (domain.OrderResponse, error) {

	var order domain.OrderResponse

	query := `SELECT * FROM orders WHERE id = $1`

	if err := i.DB.Raw(query, orderID).First(&order).Error; err != nil {
		return domain.OrderResponse{}, err
	}

	return order, nil
}

func (i *orderRepository) OrderItemsInv(productNames []string, categoryIds []int, prices, quantities []int, totalPrices []float64, userID int, orderID int) error {
	// Loop through the provided slices and perform necessary actions
	for idx := range productNames {
		orderInv := domain.OrderItemInv{
			ProductName: productNames[idx],
			Category_id: categoryIds[idx],
			Quantity:    quantities[idx],
			Price:       prices[idx],
			Total:       totalPrices[idx],
			UserID:      uint(userID),
			OrderID:     uint(orderID),
		}

		fmt.Println("order id at repo ", orderID)
		// Execute the SQL query
		query := `
			INSERT INTO order_item_invs (product_name, category_id, quantity, price, total, user_id, order_id)
			VALUES (?, ?, ?, ?, ?, ?, ?)
		`

		result := i.DB.Exec(
			query,
			orderInv.ProductName,
			orderInv.Category_id,
			orderInv.Quantity,
			orderInv.Price,
			orderInv.Total,
			orderInv.UserID,
			orderInv.OrderID,
		)

		if result.Error != nil {
			return result.Error
		}
	}

	return nil
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

func (o *orderRepository) GetOrderDetailsBrief(page int) ([]models.CombinedOrderDetails, error) {

	if page == 0 {
		page = 1
	}
	offset := (page - 1) * 3

	var orderDetails []models.CombinedOrderDetails

	err := o.DB.Raw(`
	SELECT orders.id AS order_id, orders.final_price, orders.order_status, orders.payment_status, 
	users.name, users.email, users.phone, addresses.house_name, addresses.state, 
	addresses.pin, addresses.street, addresses.city 
	FROM orders 
	INNER JOIN users ON orders.user_id = users.id 
	INNER JOIN addresses ON users.id = addresses.user_id 
	LIMIT ? OFFSET ?
`, 2, offset).Scan(&orderDetails).Error

	if err != nil {
		return []models.CombinedOrderDetails{}, nil
	}

	return orderDetails, nil
}

// CheckOrdersStatusByID retrieves the order status by ID
func (o *orderRepository) CheckOrdersStatusByID(id int) (string, error) {
	var status string
	err := o.DB.Raw("SELECT order_status FROM orders WHERE id = ?", id).Scan(&status).Error
	if err != nil {
		return "", err
	}
	return status, nil
}

// GetShipmentStatus retrieves the shipment status by order ID
func (i *orderRepository) GetShipmentStatus(orderID int) (string, error) {
	var shipmentStatus string
	err := i.DB.Exec("UPDATE orders SET order_status = 'DELIVERED', payment_status = 'PAID' WHERE id = ?", orderID).Error
	if err != nil {
		return "", err
	}
	return shipmentStatus, nil
}

func (i *orderRepository) GetOrderStatus(orderID int) (string, error) {
	var shipmentStatus string
	err := i.DB.Raw("SELECT order_status FROM orders WHERE id = ?", orderID).Scan(&shipmentStatus).Error
	if err != nil {
		return "", err
	}
	return shipmentStatus, nil
}

// ApproveOrder updates the order status to 'order_placed' for the provided order ID
func (i *orderRepository) ApproveOrder(orderID string) error {
	err := i.DB.Exec("UPDATE orders SET order_status = 'order_placed' WHERE id = ?", orderID).Error
	if err != nil {
		return err
	}
	return nil
}

// ChangeOrderStatus updates the order status for the provided order ID
func (i *orderRepository) ChangeOrderStatus(orderID int, status string) error {
	err := i.DB.Exec("UPDATE orders SET order_status = ? WHERE id = ?", status, orderID).Error
	if err != nil {
		return err
	}
	return nil
}

func (o *orderRepository) GetShipmentsStatus(orderID int) (string, error) {

	var shipmentStatus string
	err := o.DB.Raw("select order_status from orders where id = ?", orderID).Scan(&shipmentStatus).Error
	if err != nil {
		return "", err
	}

	return shipmentStatus, nil

}

func (o *orderRepository) ReturnOrder(shipmentStatus string, orderID int) error {

	err := o.DB.Exec("update orders set order_status = ?, payment_status = 'RETURNED TO WALLET' where id = ?", shipmentStatus, orderID).Error
	if err != nil {
		return err
	}

	return nil

}

func (o *orderRepository) GetOrderDetailsByOrderId(orderID string) (models.CombinedOrderDetails, error) {
	var orderDetails models.CombinedOrderDetails

	err := o.DB.Raw("SELECT orders.id, orders.final_price, orders.order_status, orders.payment_status, users.name, users.email, users.phone, addresses.house_name, addresses.state, addresses.pin, addresses.street, addresses.city "+
		"FROM orders "+
		"INNER JOIN users ON orders.user_id = users.id "+
		"INNER JOIN addresses ON users.id = addresses.user_id "+
		"WHERE orders.id = ?", orderID).Scan(&orderDetails).Error

	if err != nil {
		return models.CombinedOrderDetails{}, err
	}

	return orderDetails, nil
}

func (o *orderRepository) AddRazorPayDetails(orderID string, razorPayOrderID string) error {

	err := o.DB.Exec("insert into razer_pays (order_id,razor_id) values (?,?)", orderID, razorPayOrderID).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *orderRepository) GetOrder(orderId int) (domain.Order, error) {
	var body domain.Order
	query := `
		select * from orders
		where id = $1
	`
	if err := repo.DB.Raw(query, orderId).Scan(&body).Error; err != nil {
		return domain.Order{}, err
	}
	fmt.Println("amount", body.FinalPrice)
	return body, nil
}

func (or *orderRepository) GetOrdersDetailsByOrderId(orderID int) (models.CombinedOrderDetails, error) {

	var orderDetails models.CombinedOrderDetails
	err := or.DB.Raw(`SELECT
    orders.id as order_id,
    orders.final_price,
    orders.order_status,
    orders.payment_status,
    users.name,
    users.email,
    users.phone,
    addresses.house_name,
    addresses.state,
    addresses.street,
    addresses.city,
    addresses.pin
FROM
    orders
INNER JOIN
    users ON orders.user_id = users.id
INNER JOIN
    addresses ON users.id = addresses.user_id
WHERE
    orders.id = ?`, orderID).Scan(&orderDetails).Error
	if err != nil {
		return models.CombinedOrderDetails{}, nil
	}

	return orderDetails, nil
}

func (or *orderRepository) PaymentMethodID(orderID int) (int, error) {
	var paymentMethodID int
	err := or.DB.Raw("SELECT payment_method_id FROM orders WHERE id = ?", orderID).Scan(&paymentMethodID).Error
	if err != nil {
		return 0, err
	}
	return paymentMethodID, nil
}

func (or *orderRepository) PaymentAlreadyPaid(orderID int) (bool, error) {
	var a bool
	err := or.DB.Raw("SELECT payment_status = 'paid' FROM orders WHERE id = ?", orderID).Scan(&a).Error
	if err != nil {
		return false, err
	}
	return a, nil
}

func (repo *orderRepository) GetDetailedOrderThroughId(orderId int) (models.CombinedOrderDetails, error) {
	var body models.CombinedOrderDetails

	query := `
	SELECT 
        o.id AS order_id,
        o.final_price AS final_price,
        o.order_status AS order_status,
        o.payment_status AS payment_status,
        u.name AS name,
        u.email AS email,
        u.phone AS phone,
        a.house_name AS house_name,
        a.state AS state,
        a.pin AS pin,
        a.street AS street,
        a.city AS city
	FROM orders o
	JOIN users u ON o.user_id = u.id
	JOIN addresses a ON o.address_id = a.id 
	WHERE o.id = ?
	`
	if err := repo.DB.Raw(query, orderId).Scan(&body).Error; err != nil {
		err = errors.New("error in getting detailed order through id in repository: " + err.Error())
		return models.CombinedOrderDetails{}, err
	}
	fmt.Println("body in repo", body.OrderId)
	return body, nil
}

func (i *orderRepository) CheckPaymentStatus(orderID int) (string, error) {
	var status string

	err := i.DB.Raw("SELECT payment_status FROM orders where id = ?", orderID).Scan(&status).Error
	if err != nil {
		return "", nil
	}
	return status, err
}

func (i *orderRepository) FindFinalPrice(orderID int) (int, error) {
	var status int

	err := i.DB.Raw("SELECT final_price FROM orders where id = ?", orderID).Scan(&status).Error
	if err != nil {
		return 0, nil
	}
	return status, err
}

func (i *orderRepository) FindUserID(orderID int) (int, error) {
	var status int

	err := i.DB.Raw("SELECT user_id FROM orders where id = ?", orderID).Scan(&status).Error
	if err != nil {
		return 0, nil
	}
	return status, err
}

func (i *orderRepository) UpdateOrder(orderID int) ([]models.CombinedOrderDetails, error) {
	var body []models.CombinedOrderDetails

	err := i.DB.Exec("UPDATE orders SET order_status ='CANCELED', payment_status = 'RETURNED TO WALLET' WHERE id = ?", orderID).Error
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (i *orderRepository) UpdateReturnedOrder(orderID int) ([]models.CombinedOrderDetails, error) {
	var body []models.CombinedOrderDetails

	err := i.DB.Exec("UPDATE orders SET order_status ='RETURNED', payment_status = 'RETURNED TO WALLET' WHERE id = ?", orderID).Error
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (o *orderRepository) CheckOrderStatusByOrderId(orderID int) (string, error) {

	var status string
	err := o.DB.Raw("select order_status from orders where id = ?", orderID).Scan(&status).Error
	if err != nil {
		return "", err
	}

	return status, nil
}

func (o *orderRepository) OrderIdStatus(orderID int) (bool, error) {
	var count int
	err := o.DB.Raw("SELECT count(*) FROM orders WHERE id = ?", orderID).Row().Scan(&count)
	if err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}

func (o *orderRepository) CartExist(userID int) (bool, error) {
	var exist bool
	err := o.DB.Raw("select exists(select 1 from carts where id = ?)", userID).Scan(&exist).Error
	if err != nil {
		return false, err
	}

	return exist, nil
}

func (o *orderRepository) GetItemsByOrderId(orderId int) ([]models.ItemDetails, error) {
	var items []models.ItemDetails

	query := `
	SELECT oi.id AS order_item_id, oi.product_name, oi.quantity, oi.price, oi.total, o.id AS order_id, o.created_at, o.updated_at, o.final_price, o.order_status, o.payment_status
	FROM orders o
	JOIN order_item_invs oi ON o.id = oi.order_id
	WHERE o.id = ?;
	`

	if err := o.DB.Raw(query, orderId).Scan(&items).Error; err != nil {
		return []models.ItemDetails{}, err
	}

	return items, nil
}
