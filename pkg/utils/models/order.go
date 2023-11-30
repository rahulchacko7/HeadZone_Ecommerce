package models

type OrderDetails struct {
	OrderID         int     `json:"order_id" gorm:"column:order_id"`
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

type CombinedOrderDetails struct {
	OrderId       string  `json:"order_id"`
	FinalPrice    float64 `json:"final_price"`
	OrderStatus   string  `json:"order_status" gorm:"column:order_status"`
	PaymentStatus string  `json:"payment_status"`
	Name          string  `json:"name"`
	Email         string  `json:"email"`
	Phone         string  `json:"phone"`
	HouseName     string  `json:"house_name" validate:"required"`
	Street        string  `json:"street" validate:"required"`
	City          string  `json:"city" validate:"required"`
	State         string  `json:"state" validate:"required"`
	Pin           string  `json:"pin" validate:"required"`
}

type OrderPaymentDetails struct {
	UserID     int     `json:"user_id"`
	Username   string  `json:"username"`
	Razor_id   string  `josn:"razor_id"`
	OrderID    int     `json:"order_id"`
	FinalPrice float64 `json:"final_price"`
}
