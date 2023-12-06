package usecase

import (
	"HeadZone/pkg/repository/interfaces"
	services "HeadZone/pkg/usecase/interfaces"
	"HeadZone/pkg/utils/models"
	"errors"
)

type couponUseCase struct {
	couponRepository interfaces.CouponRepository
}

func NewCouponUseCase(repository interfaces.CouponRepository) services.CouponUseCase {
	return &couponUseCase{
		couponRepository: repository,
	}
}

func (cp *couponUseCase) AddCoupon(CouponName string, CouponStatus bool, Discount int, MinPurchase float64) (models.CouponResponse, error) {
	if Discount <= 0 || MinPurchase <= 0 {
		return models.CouponResponse{}, errors.New("discount or minimum purchase must be a +ve number")
	}
	couponResponse, err := cp.couponRepository.AddCoupon(CouponName, CouponStatus, Discount, MinPurchase)
	if err != nil {
		return models.CouponResponse{}, err
	}
	return couponResponse, nil
}

func (cp *couponUseCase) GetCoupon() ([]models.CouponResponse, error) {
	coupons, err := cp.couponRepository.GetCopupon()
	if err != nil {
		return []models.CouponResponse{}, err
	}
	return coupons, nil
}

func (cp *couponUseCase) UpdateCoupon(CId int, CouponName string, CouponStatus bool, Discount int, MinPurchase float64) (models.CouponResponse, error) {
	if Discount <= 0 || MinPurchase <= 0 || CId <= 0 {
		return models.CouponResponse{}, errors.New("discount or minimum purchase must be a +ve number")
	}
	couponResponse, err := cp.couponRepository.UpdateCoupon(CId, CouponName, CouponStatus, Discount, MinPurchase)
	if err != nil {
		return models.CouponResponse{}, err
	}
	return couponResponse, nil
}
