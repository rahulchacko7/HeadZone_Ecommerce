package usecase

import (
	"HeadZone/pkg/repository/interfaces"
	services "HeadZone/pkg/usecase/interfaces"
	"HeadZone/pkg/utils/models"
	"errors"
)

type couponUseCase struct {
	couponRepository interfaces.CouponRepository
	orderRespository interfaces.OrderRepository
	cartRepository   interfaces.CartRepository
}

func NewCouponUseCase(repository interfaces.CouponRepository, orderRepo interfaces.OrderRepository, cartRepo interfaces.CartRepository) services.CouponUseCase {
	return &couponUseCase{
		couponRepository: repository,
		orderRespository: orderRepo,
		cartRepository:   cartRepo,
	}
}

func (cp *couponUseCase) AddCoupon(CouponName string, CouponStatus bool, Discount int) (models.CouponResponse, error) {
	if Discount <= 0 {
		return models.CouponResponse{}, errors.New("discount or minimum purchase must be a +ve number")
	}
	couponResponse, err := cp.couponRepository.AddCoupon(CouponName, CouponStatus, Discount)
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

func (cp *couponUseCase) UpdateCoupon(CId int, CouponName string, CouponStatus bool, Discount int) (models.CouponResponse, error) {
	if Discount <= 0 || CId <= 0 {
		return models.CouponResponse{}, errors.New("discount or minimum purchase must be a +ve number")
	}
	couponResponse, err := cp.couponRepository.UpdateCoupon(CId, CouponName, CouponStatus, Discount)
	if err != nil {
		return models.CouponResponse{}, err
	}
	return couponResponse, nil
}

// func (repo *couponUseCase) RedeemCoupon(coupon string, UserID int) error {

// 	return errors
// }
