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
