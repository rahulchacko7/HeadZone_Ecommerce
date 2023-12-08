package models

type CouponResponse struct {
	ID                 uint    `json:"id" gorm:"primaryKey"`
	CouponName         string  `json:"coupon_name"`
	Status             bool    `json:"status" gorm:"column:status;default:true;check:status IN ('true', 'false')"`
	DiscountPercentage int     `json:"discount_percentage"`
	MinimumPrice       float64 `json:"minimum_price"`
}

type CouponDetails struct {
	CouponName string `json:"coupon_name"`
}
