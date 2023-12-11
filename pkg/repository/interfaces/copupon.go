package interfaces

import "HeadZone/pkg/utils/models"

type CouponRepository interface {
	AddCoupon(CouponName string, CouponStatus bool, Discount int) (models.CouponResponse, error)
	GetCopupon() ([]models.CouponResponse, error)
	CheckCoupon(coupon string) (bool, error)
	UpdateCoupon(CId int, CouponName string, CouponStatus bool, Discount int) (models.CouponResponse, error)
	CheckCouponValid(couponID int) (bool, error)
	FindCouponPrice(couponID int) (int, error)
}
