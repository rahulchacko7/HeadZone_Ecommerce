package models

type OrderDetails struct {
	OrderID         int     `json:"order_id" gorm:"column:id"`
	AddressID       int     `json:"address_id" gorm:"column:address_id"`
	PaymentMethodID int     `json:"payment_method_id" gorm:"column:payment_method_id"`
	Price           float64 `json:"price" gorm:"column:price"`
	OrderStatus     string  `json:"order_status" gorm:"column:order_status"`
	PaymentStatus   string  `json:"payment_status" gorm:"column:payment_status"`
}

type PaymentMethodResponse struct {
	ID           uint   `gorm:"primarykey"`
	Payment_Name string `json:"payment_name"`
}

// type AllOrderResponse struct {
// 	OrderDetails OrderDetails
// }

type Inventorories struct {
	ID          uint     `json:"id" gorm:"primaryKey"`
	CategoryID  uint     `json:"category_id"`
	Category    Category `json:"category" gorm:"foreignKey:CategoryID;constraint:OnDelete:CASCADE"`
	ProductName string   `json:"product_name"`
	Color       string   `json:"color" gorm:"color:4;default:'Black';Check:color IN ('Black', 'Blue', 'Red', 'Green');"`
	Stock       int      `json:"stock"`
	Price       float64  `json:"price"`
}
