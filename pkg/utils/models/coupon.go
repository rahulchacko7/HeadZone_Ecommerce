package models

type CouponResponse struct {
	ID           uint   `json:"id" gorm:"primaryKey"`
	CouponName   string `json:"coupon_name"`
	Status       bool   `json:"status" gorm:"column:status;default:true;check:status IN ('true', 'false')"`
	DiscountRate int    `json:"discount_rate"`
}
