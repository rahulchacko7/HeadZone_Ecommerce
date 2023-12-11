package interfaces

import "HeadZone/pkg/utils/models"

type CouponUseCase interface {
	AddCoupon(CouponName string, CouponStatus bool, Discount int) (models.CouponResponse, error)
	GetCoupon() ([]models.CouponResponse, error)
	//RedeemCoupon(coupon string, UserId int) error
	UpdateCoupon(CId int, CouponName string, CouponStatus bool, Discount int) (models.CouponResponse, error)
}
