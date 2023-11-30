package domain

type Payment struct {
	ID      uint   `json:"id" gorm:"primarykey not null"`
	OrderId int    `json:"order_id"`
	Order   Order  `json:"-" gorm:"foreignkey:OrderId"`
	RazerID string `json:"razor_id"`
	Payment string `json:"payment_id"`
}
