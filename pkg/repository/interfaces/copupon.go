package interfaces

import "HeadZone/pkg/utils/models"

type CouponRepository interface {
	AddCoupon(CouponName string, CouponStatus bool, Discount int, MinPurchase float64) (models.CouponResponse, error)
	GetCopupon() ([]models.CouponResponse, error)
	CheckCoupon(coupon string) (bool, error)
	CouponValidity(coupon string) (bool, error)
	MinimumPurchase(coupon string) (int, error)
	DiscountPercentage(coupon string) (int, error)
	UpdateUsedCoupon(coupon string, UserId int) (bool, error)
	UpdateCoupon(CId int, CouponName string, CouponStatus bool, Discount int, MinPurchase float64) (models.CouponResponse, error)
}
