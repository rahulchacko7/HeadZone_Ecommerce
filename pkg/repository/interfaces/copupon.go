package interfaces

import "HeadZone/pkg/utils/models"

type CouponRepository interface {
	AddCoupon(CouponName string, CouponStatus bool, Discount int, MinPurchase float64) (models.CouponResponse, error)
	GetCopupon() ([]models.CouponResponse, error)
	UpdateCoupon(CId int, CouponName string, CouponStatus bool, Discount int, MinPurchase float64) (models.CouponResponse, error)
}
