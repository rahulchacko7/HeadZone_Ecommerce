package models

type OrderDetails struct {
	ID            int     `json:"order_id"`
	UserName      string  `json:"name"`
	AddressID     int     `json:"address_id"`
	PaymentMethod string  `json:"payment_method"`
	Amount        float64 `json:"amount"`
}

type PaymentMethodResponse struct {
	ID           uint   `gorm:"primarykey"`
	Payment_Name string `json:"payment_name"`
}

type AllOrderResponse struct {
	OrderDetails OrderDetails
	Data         []Inventorories
}

type Inventorories struct {
	ID          uint     `json:"id" gorm:"primaryKey"`
	CategoryID  uint     `json:"category_id"`
	Category    Category `json:"category" gorm:"foreignKey:CategoryID;constraint:OnDelete:CASCADE"`
	ProductName string   `json:"product_name"`
	Color       string   `json:"color" gorm:"color:4;default:'Black';Check:color IN ('Black', 'Blue', 'Red', 'Green');"`
	Stock       int      `json:"stock"`
	Price       float64  `json:"price"`
}
